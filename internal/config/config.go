package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	RepoPath          string `mapstructure:"repo_path"`
	PreferredLanguage string `mapstructure:"preferred_language"`
	AIProvider        string `mapstructure:"ai_provider"`
	Theme             string `mapstructure:"theme"`
	TabWidth          int    `mapstructure:"tab_width"`
	TimerEnabled      bool   `mapstructure:"timer_enabled"`
}

func Load() (Config, error) {
	setDefaults()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigDir())

	viper.SetEnvPrefix("SPAR")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

func IsFirstRun() bool {
	return !viper.IsSet("repo_path")
}

func setDefaults() {
	viper.SetDefault("preferred_language", "python")
	viper.SetDefault("ai_provider", "claude")
	viper.SetDefault("theme", "dark")
	viper.SetDefault("tab_width", 4)
	viper.SetDefault("timer_enabled", true)
}
