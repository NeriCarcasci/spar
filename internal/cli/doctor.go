package cli

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/NeriCarcasci/spar/internal/ai/auth"
	"github.com/NeriCarcasci/spar/internal/config"
)

type checkResult struct {
	name   string
	ok     bool
	detail string
	remedy string
}

func RunDoctor() {
	fmt.Println(nameStyle.Render("spar"), dimStyle.Render("doctor"))
	fmt.Println()

	cfg, cfgErr := config.Load()

	var results []checkResult
	results = append(results, checkConfigFile())
	results = append(results, checkProvider(cfg, cfgErr))
	results = append(results, checkModel(cfg, cfgErr))
	results = append(results, checkAuth(cfg, cfgErr)...)
	results = append(results, checkConnectivity(cfg, cfgErr))

	issues := 0
	for _, r := range results {
		mark := successMark
		if !r.ok {
			mark = failMark
			issues++
		}
		name := fmt.Sprintf("%-22s", r.name)
		fmt.Printf("  %s %s %s\n", labelStyle.Render(name), mark, hintStyle.Render(r.detail))
		if r.remedy != "" {
			fmt.Printf("  %-22s   %s\n", "", dimStyle.Render("→ "+r.remedy))
		}
	}

	fmt.Println()
	if issues == 0 {
		PrintSuccess("All checks passed. AI features are ready.")
	} else {
		PrintError("%d issue(s) found.", issues)
	}
}

func checkConfigFile() checkResult {
	if _, err := config.Load(); err != nil {
		return checkResult{name: "Config file", ok: false, detail: "error: " + err.Error(), remedy: "Run \"spar setup\" to create configuration."}
	}
	return checkResult{name: "Config file", ok: true, detail: "found"}
}

func checkProvider(cfg config.Config, cfgErr error) checkResult {
	if cfgErr != nil {
		return checkResult{name: "AI provider", ok: false, detail: "config error"}
	}
	p := cfg.AIProvider
	if p == "" || p == "none" {
		return checkResult{
			name:   "AI provider",
			ok:     false,
			detail: "disabled",
			remedy: "Run \"spar setup\" to choose a provider.",
		}
	}
	return checkResult{name: "AI provider", ok: true, detail: p}
}

func checkModel(cfg config.Config, cfgErr error) checkResult {
	if cfgErr != nil {
		return checkResult{name: "Model", ok: false, detail: "config error"}
	}
	model := viper.GetString("ai_model")
	if model == "" {
		return checkResult{
			name:   "Model",
			ok:     false,
			detail: "not set",
			remedy: "Run \"spar settings\" to choose a model.",
		}
	}
	return checkResult{name: "Model", ok: true, detail: model}
}

func checkAuth(cfg config.Config, cfgErr error) []checkResult {
	if cfgErr != nil {
		return []checkResult{{name: "Authentication", ok: false, detail: "config error"}}
	}

	switch auth.ProviderType(cfg.AIProvider) {
	case auth.ProviderOpenAIOAuth:
		return checkOAuthAuth()
	case auth.ProviderOpenAIKey, auth.ProviderAnthropicKey, auth.ProviderOpenRouterKey:
		return checkAPIKeyAuth()
	default:
		return []checkResult{{name: "Authentication", ok: false, detail: "provider disabled", remedy: "Run \"spar setup\" to configure a provider."}}
	}
}

func checkOAuthAuth() []checkResult {
	client := auth.NewOAuthClient()

	if !client.IsAuthenticated() {
		return []checkResult{{
			name:   "Authentication",
			ok:     false,
			detail: "not authenticated",
			remedy: "Run \"spar login\" to authenticate with OpenAI.",
		}}
	}

	expires := client.TokenExpiresAt()
	if expires.IsZero() {
		return []checkResult{{name: "Authentication", ok: true, detail: "OAuth token present"}}
	}

	remaining := time.Until(expires)
	if remaining > 0 {
		return []checkResult{{
			name:   "Authentication",
			ok:     true,
			detail: fmt.Sprintf("OAuth token valid (expires in %d min)", int(remaining.Minutes())),
		}}
	}

	return []checkResult{{
		name:   "Authentication",
		ok:     false,
		detail: "OAuth token expired",
		remedy: "Token will auto-refresh on next API call, or run \"spar login\".",
	}}
}

func checkAPIKeyAuth() []checkResult {
	key := viper.GetString("ai_api_key")
	if key == "" {
		return []checkResult{{
			name:   "Authentication",
			ok:     false,
			detail: "no API key configured",
			remedy: "Run \"spar setup\" to set your API key.",
		}}
	}
	masked := key[:4] + "..." + key[len(key)-4:]
	return []checkResult{{name: "Authentication", ok: true, detail: "API key present (" + masked + ")"}}
}

func checkConnectivity(cfg config.Config, cfgErr error) checkResult {
	if cfgErr != nil {
		return checkResult{name: "API connectivity", ok: false, detail: "skipped (config error)"}
	}

	providerType := auth.ProviderType(cfg.AIProvider)

	var bearerToken string
	var baseURL string

	switch providerType {
	case auth.ProviderOpenAIOAuth:
		client := auth.NewOAuthClient()
		if !client.IsAuthenticated() {
			return checkResult{name: "API connectivity", ok: false, detail: "skipped (not authenticated)"}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		token, err := client.GetAccessToken(ctx)
		if err != nil {
			return checkResult{name: "API connectivity", ok: false, detail: "skipped (token error)"}
		}
		bearerToken = token
		baseURL = "https://api.openai.com/v1"
	case auth.ProviderOpenAIKey, auth.ProviderAnthropicKey, auth.ProviderOpenRouterKey:
		bearerToken = viper.GetString("ai_api_key")
		if bearerToken == "" {
			return checkResult{name: "API connectivity", ok: false, detail: "skipped (no API key)"}
		}
		baseURL = providerBaseURL(cfg.AIProvider)
	default:
		return checkResult{name: "API connectivity", ok: false, detail: "skipped (provider disabled)"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/models", nil)
	if err != nil {
		return checkResult{name: "API connectivity", ok: false, detail: err.Error()}
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return checkResult{name: "API connectivity", ok: false, detail: "unreachable", remedy: "Check your internet connection."}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return checkResult{name: "API connectivity", ok: true, detail: "API reachable"}
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return checkResult{
			name:   "API connectivity",
			ok:     false,
			detail: "unauthorized (401)",
			remedy: "Your credentials may be invalid. Run \"spar setup\".",
		}
	}
	return checkResult{name: "API connectivity", ok: false, detail: fmt.Sprintf("unexpected status: %d", resp.StatusCode)}
}

func providerBaseURL(provider string) string {
	switch provider {
	case "openai-key":
		return "https://api.openai.com/v1"
	case "anthropic-key":
		return "https://api.anthropic.com/v1"
	case "openrouter-key":
		return "https://openrouter.ai/api/v1"
	default:
		return ""
	}
}
