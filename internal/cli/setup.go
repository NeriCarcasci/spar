package cli

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
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
		"OpenAI (API key)",
		"Anthropic (API key — recommended)",
		"OpenRouter (API key)",
		"Skip — disable AI features",
	})

	if err := config.EnsureDirectories(); err != nil {
		PrintError("creating directories: %v", err)
		os.Exit(1)
	}

	switch choice {
	case 1:
		setupAPIKey("openai", "gpt-4o")
	case 2:
		setupAPIKey("anthropic", "claude-sonnet-4-20250514")
	case 3:
		setupAPIKey("openrouter", "gpt-4o")
	case 4:
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

func setupAPIKey(provider, defaultModel string) {
	viper.Set("ai_provider", provider)

	fmt.Println()
	model := PromptString("Choose model:", defaultModel)
	viper.Set("ai_model", model)

	fmt.Println()
	key := PromptSecret(fmt.Sprintf("Enter your %s API key:", provider))
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
