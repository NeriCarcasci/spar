package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func ConfigDir() string {
	if dir := os.Getenv("SPAR_CONFIG_DIR"); dir != "" {
		return dir
	}
	return filepath.Join(xdgConfigHome(), "spar")
}

func DataDir() string {
	if dir := os.Getenv("SPAR_DATA_DIR"); dir != "" {
		return dir
	}
	return filepath.Join(xdgDataHome(), "spar")
}

func AuthDir() string {
	return filepath.Join(ConfigDir(), "auth")
}

func ProfilePath() string {
	return filepath.Join(DataDir(), "profile.json")
}

func SessionsDir() string {
	return filepath.Join(DataDir(), "sessions")
}

func ConfigFilePath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func EnsureDirectories() error {
	dirs := []string{
		ConfigDir(),
		AuthDir(),
		DataDir(),
		SessionsDir(),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func xdgConfigHome() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir
	}
	if runtime.GOOS == "windows" {
		if dir := os.Getenv("APPDATA"); dir != "" {
			return dir
		}
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}

func xdgDataHome() string {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return dir
	}
	if runtime.GOOS == "windows" {
		if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
			return dir
		}
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share")
}
