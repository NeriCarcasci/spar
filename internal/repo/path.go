package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

func ResolveChallengesDir(repoPath string) (string, error) {
	challengesDir := filepath.Join(repoPath, "challenges")
	info, err := os.Stat(challengesDir)
	if err != nil {
		return "", fmt.Errorf("challenges directory not found at %s: %w", challengesDir, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", challengesDir)
	}
	return challengesDir, nil
}

func ResolveIndexPath(repoPath string) string {
	return filepath.Join(repoPath, "challenges", "index.yaml")
}
