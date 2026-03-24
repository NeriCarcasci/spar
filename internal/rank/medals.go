package rank

type MedalTier int

const (
	MedalNone   MedalTier = 0
	MedalBronze MedalTier = 1
	MedalSilver MedalTier = 2
	MedalGold   MedalTier = 3
)

func MedalLabel(m MedalTier) string {
	switch m {
	case MedalBronze:
		return "Bronze"
	case MedalSilver:
		return "Silver"
	case MedalGold:
		return "Gold"
	default:
		return ""
	}
}

func MedalIcon(m MedalTier) string {
	switch m {
	case MedalGold:
		return "★"
	case MedalSilver:
		return "◆"
	case MedalBronze:
		return "△"
	default:
		return "○"
	}
}

func MedalBonusSP(m MedalTier) int {
	switch m {
	case MedalBronze:
		return 200
	case MedalSilver:
		return 500
	case MedalGold:
		return 1000
	default:
		return 0
	}
}

type TrackResult struct {
	Collection       string
	TotalChallenges  int
	SolvedChallenges int
	AvgMultiplier    float64
}

func CalculateMedal(tr TrackResult) MedalTier {
	if tr.SolvedChallenges < tr.TotalChallenges {
		return MedalNone
	}
	if tr.AvgMultiplier >= 1.6 {
		return MedalGold
	}
	if tr.AvgMultiplier >= 1.3 {
		return MedalSilver
	}
	return MedalBronze
}

func UpgradeBonusSP(oldMedal, newMedal MedalTier) int {
	if newMedal <= oldMedal {
		return 0
	}
	return MedalBonusSP(newMedal) - MedalBonusSP(oldMedal)
}

var TrackNames = []string{
	"The Foundation",
	"System Design Lite",
	"Concurrency",
	"Data Structures",
	"Bit Manipulation",
	"Recursion Deep Dive",
	"Real-World Patterns",
	"Language Idiomatic",
}
