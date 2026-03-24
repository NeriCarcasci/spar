package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/NeriCarcasci/spar/internal/ai/auth"
	"github.com/NeriCarcasci/spar/internal/config"
)

func RunSetup() {
	fmt.Println(nameStyle.Render("spar"), dimStyle.Render("— setup"))
	fmt.Println()
	fmt.Println(hintStyle.Render("AI features require authentication with an LLM provider."))
	fmt.Println(hintStyle.Render("You can skip this — spar works without AI, but interview"))
	fmt.Println(hintStyle.Render("mode, hints, and post-mortem analysis won't be available."))
	fmt.Println()

	choice := PromptChoice("Choose your AI provider:", []string{
		"OpenAI (ChatGPT account — recommended, no API key needed)",
		"OpenAI (API key)",
		"Anthropic (API key)",
		"OpenRouter (API key)",
		"Skip — disable AI features",
	})

	if err := config.EnsureDirectories(); err != nil {
		PrintError("creating directories: %v", err)
		os.Exit(1)
	}

	switch choice {
	case 1:
		setupOAuth()
	case 2:
		setupAPIKey("openai-key", "gpt-4o")
	case 3:
		setupAPIKey("anthropic-key", "claude-sonnet-4-20250514")
	case 4:
		setupAPIKey("openrouter-key", "gpt-4o")
	case 5:
		viper.Set("ai_provider", "none")
		viper.Set("ai_api_key", "")
		saveConfig()
		fmt.Println()
		PrintSuccess("AI features disabled. Enable them later with \"spar settings\".")
		return
	default:
		fmt.Println("Setup cancelled.")
		return
	}

	fmt.Println()
	fmt.Println(hintStyle.Render("Run \"spar doctor\" to verify your setup."))
}

func setupOAuth() {
	viper.Set("ai_provider", string(auth.ProviderOpenAIOAuth))
	viper.Set("ai_api_key", "")

	fmt.Println()
	model := PromptString("Choose model:", "gpt-4o")
	viper.Set("ai_model", model)

	saveConfig()

	fmt.Println()
	fmt.Println(hintStyle.Render("Opening browser for OpenAI authentication..."))

	client := auth.NewOAuthClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := client.Login(ctx); err != nil {
		PrintError("Login failed: %v", err)
		fmt.Println(hintStyle.Render("You can try again with \"spar login\"."))
		return
	}

	expires := client.TokenExpiresAt()
	fmt.Println()
	PrintSuccess("Authenticated successfully")
	fmt.Printf("  %s\n", dimStyle.Render("Token expires: "+expires.Format(time.RFC822)))
}

func setupAPIKey(provider, defaultModel string) {
	viper.Set("ai_provider", provider)

	fmt.Println()
	model := PromptString("Choose model:", defaultModel)
	viper.Set("ai_model", model)

	fmt.Println()
	key := PromptSecret(fmt.Sprintf("Enter your API key:"))
	if key == "" {
		PrintError("No API key provided. Run \"spar setup\" again to configure.")
		os.Exit(1)
	}
	viper.Set("ai_api_key", key)

	saveConfig()
	fmt.Println()
	PrintSuccess("Configuration saved")
}

func saveConfig() {
	configPath := config.ConfigFilePath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := viper.SafeWriteConfigAs(configPath); err != nil {
			PrintError("saving config: %v", err)
			os.Exit(1)
		}
		return
	}
	if err := viper.WriteConfigAs(configPath); err != nil {
		PrintError("saving config: %v", err)
		os.Exit(1)
	}
}
