package rank

func StreakBonusPercent(streakDays int) float64 {
	switch {
	case streakDays >= 30:
		return 0.20
	case streakDays >= 14:
		return 0.15
	case streakDays >= 7:
		return 0.10
	case streakDays >= 3:
		return 0.05
	default:
		return 0.0
	}
}
