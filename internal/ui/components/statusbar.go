package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/ui/theme"
)

type KeyHint struct {
	Key         string
	Description string
}

type StatusBar struct {
	width int
	mode  string
	hints []KeyHint
}

func NewStatusBar() StatusBar {
	return StatusBar{}
}

func (s StatusBar) WithWidth(width int) StatusBar {
	s.width = width
	return s
}

func (s StatusBar) WithMode(mode string) StatusBar {
	s.mode = mode
	return s
}

func (s StatusBar) WithHints(hints []KeyHint) StatusBar {
	s.hints = hints
	return s
}

func (s StatusBar) View() string {
	modeStyle := lipgloss.NewStyle().
		Background(theme.Red).
		Foreground(theme.Background).
		Bold(true).
		Padding(0, 1)

	barStyle := lipgloss.NewStyle().
		Background(theme.Surface2).
		Foreground(theme.TextDim).
		Padding(0, 1)

	modeSection := modeStyle.Render(s.mode)

	var hintParts []string
	for _, h := range s.hints {
		keyStyle := lipgloss.NewStyle().
			Background(theme.Surface2).
			Foreground(theme.TextMid).
			Bold(true)
		descStyle := lipgloss.NewStyle().
			Background(theme.Surface2).
			Foreground(theme.TextDim)
		hint := keyStyle.Render(h.Key) + " " + descStyle.Render(h.Description)
		hintParts = append(hintParts, hint)
	}

	sepStyle := lipgloss.NewStyle().
		Background(theme.Surface2).
		Foreground(theme.Border)
	hintsSection := strings.Join(hintParts, sepStyle.Render(" │ "))

	modeWidth := lipgloss.Width(modeSection)
	hintsWidth := lipgloss.Width(hintsSection)
	gap := s.width - modeWidth - hintsWidth
	if gap < 0 {
		gap = 0
	}

	bar := modeSection + barStyle.Render(strings.Repeat(" ", gap)) + hintsSection

	return lipgloss.NewStyle().
		Width(s.width).
		Background(theme.Surface2).
		Render(bar)
}
