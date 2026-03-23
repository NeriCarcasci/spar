package profile

import "time"

type Profile struct {
	Username string        `json:"username"`
	Solves   []SolveRecord `json:"solves"`
}

type SolveRecord struct {
	ChallengeID string    `json:"challenge_id"`
	Language    string    `json:"language"`
	Timestamp   time.Time `json:"timestamp"`
	Duration    Duration  `json:"duration"`
	TestRuns    int       `json:"test_runs"`
	Passed      bool      `json:"passed"`
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

func (p *Profile) TotalSolved() int {
	solved := map[string]bool{}
	for _, s := range p.Solves {
		if s.Passed {
			solved[s.ChallengeID] = true
		}
	}
	return len(solved)
}

func (p *Profile) IsSolved(challengeID string) bool {
	for _, s := range p.Solves {
		if s.ChallengeID == challengeID && s.Passed {
			return true
		}
	}
	return false
}

func (p *Profile) CurrentStreak() int {
	if len(p.Solves) == 0 {
		return 0
	}

	solveDays := uniqueSolveDays(p.Solves)
	if len(solveDays) == 0 {
		return 0
	}

	today := truncateToDay(time.Now())
	lastSolve := solveDays[len(solveDays)-1]
	daysSinceLastSolve := int(today.Sub(lastSolve).Hours() / 24)
	if daysSinceLastSolve > 1 {
		return 0
	}

	streak := 1
	for i := len(solveDays) - 1; i > 0; i-- {
		diff := int(solveDays[i].Sub(solveDays[i-1]).Hours() / 24)
		if diff == 1 {
			streak++
		} else {
			break
		}
	}

	return streak
}

func (p *Profile) SolvesByDifficulty(getDifficulty func(string) string) map[string]int {
	counts := map[string]int{}
	solved := map[string]bool{}
	for _, s := range p.Solves {
		if s.Passed && !solved[s.ChallengeID] {
			solved[s.ChallengeID] = true
			difficulty := getDifficulty(s.ChallengeID)
			counts[difficulty]++
		}
	}
	return counts
}

func (p *Profile) RecentSolves(n int) []SolveRecord {
	if len(p.Solves) <= n {
		return p.Solves
	}
	return p.Solves[len(p.Solves)-n:]
}

func uniqueSolveDays(solves []SolveRecord) []time.Time {
	seen := map[time.Time]bool{}
	var days []time.Time
	for _, s := range solves {
		if !s.Passed {
			continue
		}
		day := truncateToDay(s.Timestamp)
		if !seen[day] {
			seen[day] = true
			days = append(days, day)
		}
	}
	sortDays(days)
	return days
}

func sortDays(days []time.Time) {
	for i := 1; i < len(days); i++ {
		for j := i; j > 0 && days[j].Before(days[j-1]); j-- {
			days[j], days[j-1] = days[j-1], days[j]
		}
	}
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
