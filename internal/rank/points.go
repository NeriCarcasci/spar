package rank

import "math"

type SessionInput struct {
	Difficulty    string
	DurationSecs  int
	TimeLimitSecs int
	TestRuns      int
	IsFirstSolve  bool
	UsedHints     bool
	AIEnabled     bool
	AIScore       float64
}

type SPBreakdown struct {
	Base           int
	FirstSolve     float64
	Speed          float64
	CleanRun       float64
	AIInterview    float64
	NoHint         float64
	Combined       float64
	RawSP          int
	StreakBonus    int
	StreakPercent  float64
	TotalSP        int
}

func CalculateSP(input SessionInput, streakDays int) SPBreakdown {
	bd := SPBreakdown{}

	switch input.Difficulty {
	case "easy":
		bd.Base = 10
	case "medium":
		bd.Base = 30
	case "hard":
		bd.Base = 75
	default:
		bd.Base = 10
	}

	bd.FirstSolve = 1.0
	if input.IsFirstSolve {
		bd.FirstSolve = 2.0
	}

	bd.Speed = 1.0
	if input.TimeLimitSecs > 0 && input.DurationSecs > 0 {
		ratio := float64(input.DurationSecs) / float64(input.TimeLimitSecs)
		if ratio <= 0.5 {
			bd.Speed = 1.5
		} else if ratio < 1.0 {
			bd.Speed = 1.5 - (ratio-0.5)*1.0
		}
	}

	bd.CleanRun = 1.0
	if input.TestRuns <= 1 {
		bd.CleanRun = 1.25
	}

	bd.AIInterview = 1.0
	if input.AIEnabled && input.AIScore > 0 {
		bd.AIInterview = 1.0 + input.AIScore*0.5
	}

	bd.NoHint = 1.0
	if !input.UsedHints {
		bd.NoHint = 1.1
	}

	bd.Combined = bd.FirstSolve * bd.Speed * bd.CleanRun * bd.AIInterview * bd.NoHint
	bd.RawSP = int(math.Floor(float64(bd.Base) * bd.Combined))

	bd.StreakPercent = StreakBonusPercent(streakDays)
	bd.StreakBonus = int(math.Floor(float64(bd.RawSP) * bd.StreakPercent))
	bd.TotalSP = bd.RawSP + bd.StreakBonus

	return bd
}
