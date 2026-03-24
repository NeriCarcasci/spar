package profile

import "time"

type Profile struct {
	Username        string            `json:"username"`
	TotalSP         int               `json:"total_sp"`
	CurrentTier     string            `json:"current_tier"`
	CurrentDivision int               `json:"current_division"`
	Streak          int               `json:"streak"`
	LastSolveDate   string            `json:"last_solve_date"`
	TrackMedals     map[string]int    `json:"track_medals"`
	Solves          []SolveRecord     `json:"solves"`
}

type SolveRecord struct {
	ChallengeID  string    `json:"challenge_id"`
	Language     string    `json:"language"`
	Timestamp    time.Time `json:"timestamp"`
	Duration     Duration  `json:"duration"`
	TestRuns     int       `json:"test_runs"`
	Passed       bool      `json:"passed"`
	UsedHints    bool      `json:"used_hints"`
	AIEnabled    bool      `json:"ai_enabled"`
	AIScore      float64   `json:"ai_score"`
	SPEarned     int       `json:"sp_earned"`
	Multiplier   float64   `json:"multiplier"`
	IsFirstSolve bool      `json:"is_first_solve"`
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

func (p *Profile) HasAttempted(challengeID string) bool {
	for _, s := range p.Solves {
		if s.ChallengeID == challengeID {
			return true
		}
	}
	return false
}

func (p *Profile) SPForChallenge(challengeID string) int {
	best := 0
	for _, s := range p.Solves {
		if s.ChallengeID == challengeID && s.SPEarned > best {
			best = s.SPEarned
		}
	}
	return best
}

func (p *Profile) CurrentStreak() int {
	return p.Streak
}

func (p *Profile) UpdateStreak() {
	today := time.Now().Format("2006-01-02")
	if p.LastSolveDate == today {
		return
	}
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if p.LastSolveDate == yesterday {
		p.Streak++
	} else if p.LastSolveDate == "" {
		p.Streak = 1
	} else {
		p.Streak = 1
	}
	p.LastSolveDate = today
}

func (p *Profile) IsFirstSolve(challengeID string) bool {
	for _, s := range p.Solves {
		if s.ChallengeID == challengeID && s.Passed {
			return false
		}
	}
	return true
}

func (p *Profile) EnsureDefaults() {
	if p.TrackMedals == nil {
		p.TrackMedals = make(map[string]int)
	}
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
