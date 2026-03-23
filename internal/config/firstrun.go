package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func SaveInitialConfig(repoPath string, language string, aiProvider string) error {
	if err := EnsureDirectories(); err != nil {
		return fmt.Errorf("creating config directories: %w", err)
	}

	viper.Set("repo_path", repoPath)
	viper.Set("preferred_language", language)
	viper.Set("ai_provider", aiProvider)

	configPath := ConfigFilePath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := viper.SafeWriteConfigAs(configPath); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}
		return nil
	}

	return viper.WriteConfigAs(configPath)
}
