package challenge

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ValidationError struct {
	Path    string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

func ValidateChallenge(challengePath string) []ValidationError {
	var errors []ValidationError

	challengeYAML := filepath.Join(challengePath, "challenge.yaml")
	if !fileExists(challengeYAML) {
		errors = append(errors, ValidationError{challengePath, "missing challenge.yaml"})
		return errors
	}

	testsYAML := filepath.Join(challengePath, "tests.yaml")
	if !fileExists(testsYAML) {
		errors = append(errors, ValidationError{challengePath, "missing tests.yaml"})
	}

	setupDir := filepath.Join(challengePath, "setup")
	if !dirExists(setupDir) {
		errors = append(errors, ValidationError{challengePath, "missing setup/ directory"})
	}

	data, err := os.ReadFile(challengeYAML)
	if err != nil {
		errors = append(errors, ValidationError{challengeYAML, fmt.Sprintf("cannot read: %v", err)})
		return errors
	}

	errors = append(errors, validateChallengeContent(challengePath, data)...)

	return errors
}

func validateChallengeContent(challengePath string, data []byte) []ValidationError {
	var errors []ValidationError
	var ch Challenge

	if err := yamlUnmarshal(data, &ch); err != nil {
		errors = append(errors, ValidationError{challengePath, fmt.Sprintf("invalid YAML: %v", err)})
		return errors
	}

	if ch.ID == "" {
		errors = append(errors, ValidationError{challengePath, "missing id field"})
	}
	if ch.Title == "" {
		errors = append(errors, ValidationError{challengePath, "missing title field"})
	}
	if ch.Difficulty == "" {
		errors = append(errors, ValidationError{challengePath, "missing difficulty field"})
	}
	if ch.Category == "" {
		errors = append(errors, ValidationError{challengePath, "missing category field"})
	}
	if len(ch.Languages) == 0 {
		errors = append(errors, ValidationError{challengePath, "no languages declared"})
	}
	if ch.Description == "" {
		errors = append(errors, ValidationError{challengePath, "missing description"})
	}

	if !isValidDifficulty(ch.Difficulty) {
		errors = append(errors, ValidationError{challengePath, fmt.Sprintf("invalid difficulty: %s", ch.Difficulty)})
	}

	for _, lang := range ch.Languages {
		setupFile := filepath.Join(challengePath, "setup", setupFilename(lang))
		if !fileExists(setupFile) {
			errors = append(errors, ValidationError{challengePath, fmt.Sprintf("missing setup file for %s", lang)})
		}
	}

	return errors
}

func ValidateAll(challengesDir string) []ValidationError {
	var allErrors []ValidationError

	categories, err := os.ReadDir(challengesDir)
	if err != nil {
		return []ValidationError{{challengesDir, fmt.Sprintf("cannot read directory: %v", err)}}
	}

	for _, category := range categories {
		if !category.IsDir() {
			continue
		}
		categoryPath := filepath.Join(challengesDir, category.Name())
		challenges, err := os.ReadDir(categoryPath)
		if err != nil {
			allErrors = append(allErrors, ValidationError{categoryPath, fmt.Sprintf("cannot read: %v", err)})
			continue
		}
		for _, ch := range challenges {
			if !ch.IsDir() {
				continue
			}
			challengePath := filepath.Join(categoryPath, ch.Name())
			allErrors = append(allErrors, ValidateChallenge(challengePath)...)
		}
	}

	return allErrors
}

func isValidDifficulty(d Difficulty) bool {
	return d == Easy || d == Medium || d == Hard
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func yamlUnmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
