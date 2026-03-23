package challenge

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadIndex(repoPath string) (*Index, error) {
	indexPath := filepath.Join(repoPath, "challenges", "index.yaml")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("reading index at %s: %w", indexPath, err)
	}

	var index Index
	if err := yaml.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("parsing index: %w", err)
	}

	return &index, nil
}

func (idx *Index) FindByID(id string) (IndexEntry, bool) {
	for _, entry := range idx.Challenges {
		if entry.ID == id {
			return entry, true
		}
	}
	return IndexEntry{}, false
}

func (idx *Index) FilterByCategory(category string) []IndexEntry {
	var results []IndexEntry
	for _, entry := range idx.Challenges {
		if entry.Category == category {
			results = append(results, entry)
		}
	}
	return results
}

func (idx *Index) FilterByDifficulty(difficulty Difficulty) []IndexEntry {
	var results []IndexEntry
	for _, entry := range idx.Challenges {
		if entry.Difficulty == difficulty {
			results = append(results, entry)
		}
	}
	return results
}

func (idx *Index) FilterByLanguage(language string) []IndexEntry {
	var results []IndexEntry
	for _, entry := range idx.Challenges {
		if containsString(entry.Languages, language) {
			results = append(results, entry)
		}
	}
	return results
}

func (idx *Index) Categories() []string {
	seen := map[string]bool{}
	var categories []string
	for _, entry := range idx.Challenges {
		if !seen[entry.Category] {
			seen[entry.Category] = true
			categories = append(categories, entry.Category)
		}
	}
	return categories
}

func containsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
