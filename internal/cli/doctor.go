package cli

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
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
	results = append(results, checkAPIKey(cfg, cfgErr))
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
	path := config.ConfigFilePath()
	if _, err := config.Load(); err != nil {
		return checkResult{name: "Config file", ok: false, detail: "error: " + err.Error(), remedy: "Run \"spar setup\" to create configuration."}
	}
	_ = path
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

func checkAPIKey(cfg config.Config, cfgErr error) checkResult {
	if cfgErr != nil {
		return checkResult{name: "API key", ok: false, detail: "config error"}
	}
	p := cfg.AIProvider
	if p == "" || p == "none" {
		return checkResult{name: "API key", ok: false, detail: "skipped (no provider)"}
	}

	key := viper.GetString("ai_api_key")
	if key == "" {
		return checkResult{
			name:   "API key",
			ok:     false,
			detail: "not configured",
			remedy: "Run \"spar setup\" to set your API key.",
		}
	}

	masked := key[:4] + "..." + key[len(key)-4:]
	return checkResult{name: "API key", ok: true, detail: "present (" + masked + ")"}
}

func checkConnectivity(cfg config.Config, cfgErr error) checkResult {
	if cfgErr != nil {
		return checkResult{name: "API connectivity", ok: false, detail: "skipped (config error)"}
	}

	key := viper.GetString("ai_api_key")
	if key == "" {
		return checkResult{name: "API connectivity", ok: false, detail: "skipped (no API key)"}
	}

	baseURL := providerBaseURL(cfg.AIProvider)
	if baseURL == "" {
		return checkResult{name: "API connectivity", ok: false, detail: "unknown provider"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/models", nil)
	if err != nil {
		return checkResult{name: "API connectivity", ok: false, detail: err.Error()}
	}
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return checkResult{
			name:   "API connectivity",
			ok:     false,
			detail: "unreachable",
			remedy: "Check your internet connection.",
		}
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
			remedy: "Your API key may be invalid. Run \"spar setup\" to reconfigure.",
		}
	}
	return checkResult{
		name:   "API connectivity",
		ok:     false,
		detail: fmt.Sprintf("unexpected status: %d", resp.StatusCode),
	}
}

func providerBaseURL(provider string) string {
	switch provider {
	case "openai":
		return "https://api.openai.com/v1"
	case "anthropic":
		return "https://api.anthropic.com/v1"
	case "openrouter":
		return "https://openrouter.ai/api/v1"
	default:
		return ""
	}
}
