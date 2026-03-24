package rank

type Tier struct {
	Name      string
	Color     string
	Thresholds [3]int
}

type Division int

const (
	DivisionIII Division = 0
	DivisionII  Division = 1
	DivisionI   Division = 2
)

type RankInfo struct {
	Tier       Tier
	TierIndex  int
	Division   Division
	TotalSP    int
	Progress   float64
	NextSP     int
	NextName   string
	IsMax      bool
}

var Tiers = []Tier{
	{Name: "Spark", Color: "#555555", Thresholds: [3]int{0, 50, 150}},
	{Name: "Cipher", Color: "#5DCAA5", Thresholds: [3]int{300, 500, 750}},
	{Name: "Warden", Color: "#97C459", Thresholds: [3]int{1100, 1500, 2000}},
	{Name: "Sentinel", Color: "#FBBF24", Thresholds: [3]int{2700, 3500, 4500}},
	{Name: "Arbiter", Color: "#F0997B", Thresholds: [3]int{5800, 7500, 9500}},
	{Name: "Sovereign", Color: "#FF3B30", Thresholds: [3]int{12000, 15000, 19000}},
	{Name: "Mythic", Color: "#FF3B30", Thresholds: [3]int{24000, 30000, 40000}},
}

func DivisionLabel(d Division) string {
	switch d {
	case DivisionIII:
		return "III"
	case DivisionII:
		return "II"
	case DivisionI:
		return "I"
	default:
		return ""
	}
}

func Calculate(totalSP int) RankInfo {
	info := RankInfo{TotalSP: totalSP}

	for ti := len(Tiers) - 1; ti >= 0; ti-- {
		for di := DivisionI; di >= DivisionIII; di-- {
			threshold := Tiers[ti].Thresholds[di]
			if totalSP >= threshold {
				info.Tier = Tiers[ti]
				info.TierIndex = ti
				info.Division = di

				nextTier, nextDiv := nextStep(ti, int(di))
				if nextTier < 0 {
					info.IsMax = true
					info.Progress = 1.0
					info.NextSP = threshold
					info.NextName = Tiers[ti].Name + " " + DivisionLabel(di)
				} else {
					nextThreshold := Tiers[nextTier].Thresholds[nextDiv]
					info.NextSP = nextThreshold
					info.NextName = Tiers[nextTier].Name + " " + DivisionLabel(Division(nextDiv))
					span := nextThreshold - threshold
					if span > 0 {
						info.Progress = float64(totalSP-threshold) / float64(span)
					}
				}

				return info
			}
		}
	}

	info.Tier = Tiers[0]
	info.TierIndex = 0
	info.Division = DivisionIII
	info.NextSP = Tiers[0].Thresholds[1]
	info.NextName = Tiers[0].Name + " " + DivisionLabel(DivisionII)
	if info.NextSP > 0 {
		info.Progress = float64(totalSP) / float64(info.NextSP)
	}
	return info
}

func nextStep(tierIdx, divIdx int) (int, int) {
	if divIdx < 2 {
		return tierIdx, divIdx + 1
	}
	if tierIdx < len(Tiers)-1 {
		return tierIdx + 1, 0
	}
	return -1, -1
}

func FullName(ri RankInfo) string {
	return ri.Tier.Name + " " + DivisionLabel(ri.Division)
}
