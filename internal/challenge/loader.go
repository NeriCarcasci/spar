package challenge

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadChallenge(repoPath string, entry IndexEntry) (*Challenge, error) {
	challengePath := filepath.Join(repoPath, entry.Path, "challenge.yaml")
	data, err := os.ReadFile(challengePath)
	if err != nil {
		return nil, fmt.Errorf("reading challenge %s: %w", entry.ID, err)
	}

	var ch Challenge
	if err := yaml.Unmarshal(data, &ch); err != nil {
		return nil, fmt.Errorf("parsing challenge %s: %w", entry.ID, err)
	}

	ch.Path = filepath.Join(repoPath, entry.Path)
	return &ch, nil
}

func LoadTests(challengePath string) (*TestSuite, error) {
	testsPath := filepath.Join(challengePath, "tests.yaml")
	data, err := os.ReadFile(testsPath)
	if err != nil {
		return nil, fmt.Errorf("reading tests at %s: %w", testsPath, err)
	}

	var suite TestSuite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("parsing tests: %w", err)
	}

	return &suite, nil
}

func LoadSetupCode(challengePath string, language string) (string, error) {
	filename := setupFilename(language)
	setupPath := filepath.Join(challengePath, "setup", filename)
	data, err := os.ReadFile(setupPath)
	if err != nil {
		return "", fmt.Errorf("reading setup for %s: %w", language, err)
	}
	return string(data), nil
}

func setupFilename(language string) string {
	switch language {
	case "python":
		return "python.py"
	case "go":
		return "solution.go"
	case "javascript":
		return "javascript.js"
	case "cpp":
		return "cpp.cpp"
	case "rust":
		return "rust.rs"
	default:
		return language + ".txt"
	}
}
