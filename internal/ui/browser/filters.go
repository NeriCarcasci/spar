package browser

import "github.com/spar-cli/spar/internal/challenge"

type Filters struct {
	Category     string
	Difficulty   challenge.Difficulty
	Language     string
	SolvedOnly   bool
	UnsolvedOnly bool
}

func (f Filters) IsActive() bool {
	return f.Category != "" || f.Difficulty != "" || f.Language != "" || f.SolvedOnly || f.UnsolvedOnly
}

func (f Filters) Apply(entries []challenge.IndexEntry, isSolved func(string) bool) []challenge.IndexEntry {
	var results []challenge.IndexEntry
	for _, entry := range entries {
		if f.matches(entry, isSolved) {
			results = append(results, entry)
		}
	}
	return results
}

func (f Filters) matches(entry challenge.IndexEntry, isSolved func(string) bool) bool {
	if f.Category != "" && entry.Category != f.Category {
		return false
	}
	if f.Difficulty != "" && entry.Difficulty != f.Difficulty {
		return false
	}
	if f.Language != "" && !containsLanguage(entry.Languages, f.Language) {
		return false
	}
	if f.SolvedOnly && !isSolved(entry.ID) {
		return false
	}
	if f.UnsolvedOnly && isSolved(entry.ID) {
		return false
	}
	return true
}

func containsLanguage(languages []string, target string) bool {
	for _, lang := range languages {
		if lang == target {
			return true
		}
	}
	return false
}
