package auth

import (
	"context"
	"fmt"
	"net/http"
)

type ProviderType string

const (
	ProviderOpenAIOAuth   ProviderType = "openai-oauth"
	ProviderOpenAIKey     ProviderType = "openai-key"
	ProviderAnthropicKey  ProviderType = "anthropic-key"
	ProviderOpenRouterKey ProviderType = "openrouter-key"
	ProviderNone          ProviderType = "none"
)

type AuthProvider interface {
	GetAuthHeader(ctx context.Context) (string, error)
	BaseURL() string
	IsAuthenticated() bool
	ProviderName() ProviderType
}

type OAuthAuthProvider struct {
	client *OAuthClient
}

func NewOAuthAuthProvider() *OAuthAuthProvider {
	return &OAuthAuthProvider{client: NewOAuthClient()}
}

func (p *OAuthAuthProvider) GetAuthHeader(ctx context.Context) (string, error) {
	token, err := p.client.GetAccessToken(ctx)
	if err != nil {
		return "", err
	}
	return "Bearer " + token, nil
}

func (p *OAuthAuthProvider) BaseURL() string {
	return "https://api.openai.com/v1"
}

func (p *OAuthAuthProvider) IsAuthenticated() bool {
	return p.client.IsAuthenticated()
}

func (p *OAuthAuthProvider) ProviderName() ProviderType {
	return ProviderOpenAIOAuth
}

func (p *OAuthAuthProvider) OAuthClient() *OAuthClient {
	return p.client
}

type APIKeyAuthProvider struct {
	providerType ProviderType
	apiKey       string
	baseURL      string
}

func NewAPIKeyAuthProvider(providerType ProviderType, apiKey, baseURL string) *APIKeyAuthProvider {
	if baseURL == "" {
		switch providerType {
		case ProviderOpenAIKey:
			baseURL = "https://api.openai.com/v1"
		case ProviderAnthropicKey:
			baseURL = "https://api.anthropic.com/v1"
		case ProviderOpenRouterKey:
			baseURL = "https://openrouter.ai/api/v1"
		}
	}
	return &APIKeyAuthProvider{
		providerType: providerType,
		apiKey:       apiKey,
		baseURL:      baseURL,
	}
}

func (p *APIKeyAuthProvider) GetAuthHeader(ctx context.Context) (string, error) {
	if p.apiKey == "" {
		return "", fmt.Errorf("no API key configured for %s", p.providerType)
	}
	return "Bearer " + p.apiKey, nil
}

func (p *APIKeyAuthProvider) BaseURL() string {
	return p.baseURL
}

func (p *APIKeyAuthProvider) IsAuthenticated() bool {
	return p.apiKey != ""
}

func (p *APIKeyAuthProvider) ProviderName() ProviderType {
	return p.providerType
}

type NoopAuthProvider struct{}

func (p *NoopAuthProvider) GetAuthHeader(ctx context.Context) (string, error) {
	return "", fmt.Errorf("AI features disabled — no provider configured")
}

func (p *NoopAuthProvider) BaseURL() string {
	return ""
}

func (p *NoopAuthProvider) IsAuthenticated() bool {
	return false
}

func (p *NoopAuthProvider) ProviderName() ProviderType {
	return ProviderNone
}

func ApplyAuth(ctx context.Context, provider AuthProvider, req *http.Request) error {
	header, err := provider.GetAuthHeader(ctx)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", header)
	return nil
}
