package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/NeriCarcasci/spar/internal/challenge"
	"gopkg.in/yaml.v3"
)

func GenerateIndex(challengesDir string) error {
	entries, err := walkChallenges(challengesDir)
	if err != nil {
		return fmt.Errorf("walking challenges: %w", err)
	}

	sortEntries(entries)

	index := challenge.Index{Challenges: entries}
	data, err := yaml.Marshal(index)
	if err != nil {
		return fmt.Errorf("marshaling index: %w", err)
	}

	outputPath := filepath.Join(challengesDir, "index.yaml")
	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return fmt.Errorf("writing index: %w", err)
	}

	return nil
}

func walkChallenges(challengesDir string) ([]challenge.IndexEntry, error) {
	var entries []challenge.IndexEntry

	categories, err := os.ReadDir(challengesDir)
	if err != nil {
		return nil, fmt.Errorf("reading challenges directory: %w", err)
	}

	for _, category := range categories {
		if !category.IsDir() {
			continue
		}
		categoryPath := filepath.Join(challengesDir, category.Name())
		challenges, err := os.ReadDir(categoryPath)
		if err != nil {
			continue
		}

		for _, ch := range challenges {
			if !ch.IsDir() {
				continue
			}
			entry, err := readChallengeEntry(challengesDir, category.Name(), ch.Name())
			if err != nil {
				continue
			}
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func readChallengeEntry(challengesDir string, categoryName string, challengeName string) (challenge.IndexEntry, error) {
	challengePath := filepath.Join(challengesDir, categoryName, challengeName)
	yamlPath := filepath.Join(challengePath, "challenge.yaml")

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return challenge.IndexEntry{}, fmt.Errorf("reading %s: %w", yamlPath, err)
	}

	var ch challenge.Challenge
	if err := yaml.Unmarshal(data, &ch); err != nil {
		return challenge.IndexEntry{}, fmt.Errorf("parsing %s: %w", yamlPath, err)
	}

	relativePath := filepath.Join("challenges", categoryName, challengeName)

	return challenge.IndexEntry{
		ID:         ch.ID,
		Title:      ch.Title,
		Difficulty: ch.Difficulty,
		Category:   ch.Category,
		Tags:       ch.Tags,
		Languages:  ch.Languages,
		Path:       filepath.ToSlash(relativePath),
	}, nil
}

func sortEntries(entries []challenge.IndexEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Category != entries[j].Category {
			return entries[i].Category < entries[j].Category
		}
		if entries[i].Difficulty != entries[j].Difficulty {
			return difficultyRank(entries[i].Difficulty) < difficultyRank(entries[j].Difficulty)
		}
		return entries[i].Title < entries[j].Title
	})
}

func difficultyRank(d challenge.Difficulty) int {
	switch d {
	case challenge.Easy:
		return 0
	case challenge.Medium:
		return 1
	case challenge.Hard:
		return 2
	default:
		return 3
	}
}
