package cli

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/NeriCarcasci/spar/internal/config"
)

func RunSettings() {
	cfg, err := config.Load()
	if err != nil {
		PrintError("loading config: %v", err)
		os.Exit(1)
	}

	fmt.Println(nameStyle.Render("spar"), dimStyle.Render("settings"))
	fmt.Println()
	fmt.Printf("  %s %s\n", labelStyle.Render("Provider:          "), hintStyle.Render(cfg.AIProvider))
	fmt.Printf("  %s %s\n", labelStyle.Render("Model:             "), hintStyle.Render(viper.GetString("ai_model")))
	fmt.Printf("  %s %s\n", labelStyle.Render("Language:          "), hintStyle.Render(cfg.PreferredLanguage))

	key := viper.GetString("ai_api_key")
	if key != "" && len(key) >= 8 {
		masked := key[:4] + "..." + key[len(key)-4:]
		fmt.Printf("  %s %s\n", labelStyle.Render("API key:           "), dimStyle.Render(masked))
	}
	fmt.Println()

	for {
		choice := PromptChoice("What would you like to change?", []string{
			"AI provider + API key",
			"Model",
			"Preferred language",
			"Done",
		})

		switch choice {
		case 1:
			changeProvider()
		case 2:
			changeModel()
		case 3:
			changeLanguage()
		case 4:
			return
		default:
			return
		}

		saveConfig()
		PrintSuccess("Settings saved")
		fmt.Println()
	}
}

func RunSettingsReset() {
	if !PromptYesNo("This will delete all spar configuration and stored data. Continue?", false) {
		fmt.Println("Cancelled.")
		return
	}

	configPath := config.ConfigFilePath()
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		PrintError("removing config: %v", err)
		os.Exit(1)
	}

	authDir := config.AuthDir()
	if err := os.RemoveAll(authDir); err != nil {
		PrintError("removing auth data: %v", err)
		os.Exit(1)
	}

	PrintSuccess("Configuration and credentials removed.")
	fmt.Println(hintStyle.Render("Run \"spar setup\" to reconfigure."))
}

func changeProvider() {
	choice := PromptChoice("Choose your AI provider:", []string{
		"OpenAI (API key)",
		"Anthropic (API key)",
		"OpenRouter (API key)",
		"Disable AI features",
	})

	switch choice {
	case 1:
		viper.Set("ai_provider", "openai-key")
		fmt.Println()
		key := PromptSecret("Enter your OpenAI API key:")
		viper.Set("ai_api_key", key)
	case 2:
		viper.Set("ai_provider", "anthropic-key")
		fmt.Println()
		key := PromptSecret("Enter your Anthropic API key:")
		viper.Set("ai_api_key", key)
	case 3:
		viper.Set("ai_provider", "openrouter-key")
		fmt.Println()
		key := PromptSecret("Enter your OpenRouter API key:")
		viper.Set("ai_api_key", key)
	case 4:
		viper.Set("ai_provider", "none")
		viper.Set("ai_api_key", "")
	}
}

func changeModel() {
	current := viper.GetString("ai_model")
	model := PromptString("Choose model:", current)
	viper.Set("ai_model", model)
}

func changeLanguage() {
	choice := PromptChoice("Choose your preferred language:", []string{
		"Python",
		"Go",
		"JavaScript",
		"TypeScript",
	})
	langs := []string{"python", "go", "javascript", "typescript"}
	if choice >= 1 && choice <= len(langs) {
		viper.Set("preferred_language", langs[choice-1])
	}
}
