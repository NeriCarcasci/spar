package rank

import (
	"github.com/charmbracelet/lipgloss"
)

func RenderIcon(tierIndex int) string {
	if tierIndex < 0 || tierIndex >= len(Tiers) {
		tierIndex = 0
	}
	color := lipgloss.Color(Tiers[tierIndex].Color)
	s := lipgloss.NewStyle().Foreground(color)

	switch tierIndex {
	case 0:
		return s.Render("  ·  \n ·∘· \n  ·  ")
	case 1:
		return s.Render(" ╭─╮ \n │◆│ \n ╰─╯ ")
	case 2:
		return s.Render("  △  \n ╔═╗ \n ╚═╝ ")
	case 3:
		return s.Render("  ◇  \n ◁◆▷ \n  │  ")
	case 4:
		return s.Render(" ╱◆╲ \n ╲▪╱ \n  ═  ")
	case 5:
		return s.Render(" ╲★╱ \n ╱◆╲ \n ═╪═ ")
	case 6:
		outer := lipgloss.NewStyle().Foreground(color)
		accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B5E"))
		return outer.Render(" ╋") + accent.Render("★") + outer.Render("╋ \n ╲") + accent.Render("◆") + outer.Render("╱ \n ═") + accent.Render("╬") + outer.Render("═ ")
	default:
		return s.Render("  ·  \n ·∘· \n  ·  ")
	}
}

func RenderInline(tierIndex int) string {
	if tierIndex < 0 || tierIndex >= len(Tiers) {
		tierIndex = 0
	}
	color := lipgloss.Color(Tiers[tierIndex].Color)
	s := lipgloss.NewStyle().Foreground(color)

	switch tierIndex {
	case 0:
		return s.Render("∘")
	case 1:
		return s.Render("◆")
	case 2:
		return s.Render("△")
	case 3:
		return s.Render("◁◆▷")
	case 4:
		return s.Render("◆")
	case 5:
		return s.Render("★")
	case 6:
		accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B5E"))
		return accent.Render("★")
	default:
		return s.Render("∘")
	}
}
