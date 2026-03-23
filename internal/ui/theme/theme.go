package theme

import "github.com/charmbracelet/lipgloss"

var (
	Background = lipgloss.Color("#0A0A0A")
	Surface    = lipgloss.Color("#111111")
	Surface2   = lipgloss.Color("#1A1A1A")
	Surface3   = lipgloss.Color("#222222")
	Border     = lipgloss.Color("#2A2A2A")

	TextPrimary = lipgloss.Color("#E8E8E8")
	TextMid     = lipgloss.Color("#999999")
	TextDim     = lipgloss.Color("#555555")

	Red      = lipgloss.Color("#FF3B30")
	RedLight = lipgloss.Color("#FF6B5E")
	RedDim   = lipgloss.Color("#CC2F26")

	Green    = lipgloss.Color("#4ADE80")
	GreenDim = lipgloss.Color("#1A2E1A")

	Amber    = lipgloss.Color("#FBBF24")
	AmberDim = lipgloss.Color("#2E2A1A")

	Purple    = lipgloss.Color("#A78BFA")
	PurpleDim = lipgloss.Color("#1E1A2E")

	HardBg = lipgloss.Color("#2E1A1A")
)

func DifficultyStyle(difficulty string) lipgloss.Style {
	switch difficulty {
	case "easy":
		return lipgloss.NewStyle().
			Foreground(Green).
			Background(GreenDim).
			Padding(0, 1).
			Bold(true)
	case "medium":
		return lipgloss.NewStyle().
			Foreground(Amber).
			Background(AmberDim).
			Padding(0, 1).
			Bold(true)
	case "hard":
		return lipgloss.NewStyle().
			Foreground(Red).
			Background(HardBg).
			Padding(0, 1).
			Bold(true)
	default:
		return lipgloss.NewStyle().
			Foreground(TextDim).
			Padding(0, 1)
	}
}

func CategoryBadge() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(Purple).
		Background(PurpleDim).
		Padding(0, 1)
}

func LanguageBadge() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(TextMid).
		Background(Surface2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(Border).
		Padding(0, 1)
}
