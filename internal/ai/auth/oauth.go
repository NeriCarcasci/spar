package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/NeriCarcasci/spar/internal/config"
)

const (
	clientID            = "app_EMoamEEZ73f0CkXaXp7hrann"
	authorizationURL    = "https://auth.openai.com/oauth/authorize"
	tokenURL            = "https://auth.openai.com/oauth/token"
	scope               = "openai.chat"
	tokenRefreshMargin  = 5 * time.Minute
	callbackPath        = "/auth/callback"
	portRangeStart      = 14550
	portRangeEnd        = 14600
	httpClientTimeout   = 30 * time.Second
	serverShutdownTimout = 5 * time.Second
)

type StoredTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	Provider     string `json:"provider"`
}

type OAuthClient struct {
	mu     sync.Mutex
	tokens *StoredTokens
	http   *http.Client
}

func NewOAuthClient() *OAuthClient {
	return &OAuthClient{
		http: &http.Client{Timeout: httpClientTimeout},
	}
}

func (o *OAuthClient) Login(ctx context.Context) error {
	verifier, challenge, err := generatePKCE()
	if err != nil {
		return fmt.Errorf("generate PKCE: %w", err)
	}

	state, err := generateRandomString(32)
	if err != nil {
		return fmt.Errorf("generate state: %w", err)
	}

	listener, port, err := findAvailablePort()
	if err != nil {
		return fmt.Errorf("find port: %w", err)
	}

	redirectURI := fmt.Sprintf("http://localhost:%d%s", port, callbackPath)
	authURL := buildAuthorizationURL(redirectURI, challenge, state)

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	server := &http.Server{Handler: mux}

	mux.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			errCh <- fmt.Errorf("state mismatch")
			http.Error(w, "State mismatch", http.StatusBadRequest)
			return
		}

		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			desc := r.URL.Query().Get("error_description")
			errCh <- fmt.Errorf("oauth error: %s — %s", errMsg, desc)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in callback")
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, callbackHTML)
		codeCh <- code
	})

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("callback server: %w", err)
		}
	}()

	defer func() {
		shutCtx, cancel := context.WithTimeout(context.Background(), serverShutdownTimout)
		defer cancel()
		server.Shutdown(shutCtx)
	}()

	if err := openBrowser(authURL); err != nil {
		return fmt.Errorf("open browser: %w", err)
	}

	fmt.Printf("Waiting for callback on http://localhost:%d%s\n", port, callbackPath)

	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Minute):
		return fmt.Errorf("login timed out waiting for callback")
	}

	tokens, err := o.exchangeCode(ctx, code, redirectURI, verifier)
	if err != nil {
		return fmt.Errorf("exchange code: %w", err)
	}

	o.mu.Lock()
	o.tokens = tokens
	o.mu.Unlock()

	if err := storeTokens(tokens); err != nil {
		return fmt.Errorf("store tokens: %w", err)
	}

	return nil
}

func (o *OAuthClient) Logout() error {
	o.mu.Lock()
	o.tokens = nil
	o.mu.Unlock()
	return deleteTokens()
}

func (o *OAuthClient) GetAccessToken(ctx context.Context) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.tokens == nil {
		loaded, err := loadTokens()
		if err != nil {
			return "", fmt.Errorf("not authenticated: %w", err)
		}
		o.tokens = loaded
	}

	if time.Now().Unix() >= o.tokens.ExpiresAt-int64(tokenRefreshMargin.Seconds()) {
		refreshed, err := o.refreshToken(ctx, o.tokens.RefreshToken)
		if err != nil {
			o.tokens = nil
			return "", fmt.Errorf("token refresh failed (re-login required): %w", err)
		}
		o.tokens = refreshed
		if err := storeTokens(refreshed); err != nil {
			return "", fmt.Errorf("store refreshed tokens: %w", err)
		}
	}

	return o.tokens.AccessToken, nil
}

func (o *OAuthClient) IsAuthenticated() bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.tokens != nil {
		return true
	}
	loaded, err := loadTokens()
	if err != nil {
		return false
	}
	o.tokens = loaded
	return true
}

func (o *OAuthClient) TokenExpiresAt() time.Time {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.tokens == nil {
		return time.Time{}
	}
	return time.Unix(o.tokens.ExpiresAt, 0)
}

func (o *OAuthClient) exchangeCode(ctx context.Context, code, redirectURI, verifier string) (*StoredTokens, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"client_id":     {clientID},
		"code_verifier": {verifier},
	}
	return o.tokenRequest(ctx, data)
}

func (o *OAuthClient) refreshToken(ctx context.Context, refreshToken string) (*StoredTokens, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {clientID},
	}
	return o.tokenRequest(ctx, data)
}

func (o *OAuthClient) tokenRequest(ctx context.Context, data url.Values) (*StoredTokens, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("token endpoint returned %d: %s — %s", resp.StatusCode, errResp.Error, errResp.Description)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}

	return &StoredTokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Unix() + tokenResp.ExpiresIn,
		Provider:     "openai",
	}, nil
}

func tokenFilePath() string {
	return filepath.Join(config.AuthDir(), "openai.json")
}

func storeTokens(tokens *StoredTokens) error {
	path := tokenFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func loadTokens() (*StoredTokens, error) {
	data, err := os.ReadFile(tokenFilePath())
	if err != nil {
		return nil, err
	}
	var tokens StoredTokens
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, err
	}
	return &tokens, nil
}

func deleteTokens() error {
	err := os.Remove(tokenFilePath())
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func generatePKCE() (verifier string, challenge string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", err
	}
	verifier = base64.RawURLEncoding.EncodeToString(buf)
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	return verifier, challenge, nil
}

func generateRandomString(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func buildAuthorizationURL(redirectURI, challenge, state string) string {
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {clientID},
		"redirect_uri":          {redirectURI},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
		"state":                 {state},
		"scope":                 {scope},
	}
	return authorizationURL + "?" + params.Encode()
}

func findAvailablePort() (net.Listener, int, error) {
	for port := portRangeStart; port <= portRangeEnd; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			return ln, port, nil
		}
	}
	return nil, 0, fmt.Errorf("no available port in range %d-%d", portRangeStart, portRangeEnd)
}

func openBrowser(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "linux":
		cmd = exec.Command("xdg-open", rawURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", rawURL)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}

const callbackHTML = `<!DOCTYPE html>
<html><head><title>spar</title>
<style>body{font-family:system-ui,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0;background:#0a0a0a;color:#e0e0e0}
.card{text-align:center;padding:2rem;border-radius:12px;background:#1a1a1a;border:1px solid #333}
h1{color:#4ade80;margin:0 0 .5rem}p{color:#999;margin:0}</style></head>
<body><div class="card"><h1>Authenticated</h1><p>Return to spar. You can close this tab.</p></div></body></html>`
