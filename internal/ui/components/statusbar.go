package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/NeriCarcasci/spar/internal/ui/theme"
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
	if s.width <= 0 {
		return ""
	}

	modeStyle := lipgloss.NewStyle().
		Background(theme.Red).
		Foreground(theme.Background).
		Bold(true).
		Padding(0, 1)

	hintKeyStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		Bold(true)

	hintTextStyle := lipgloss.NewStyle().
		Foreground(theme.TextFaint)

	separatorStyle := lipgloss.NewStyle().Foreground(theme.TextFaint)
	barStyle := lipgloss.NewStyle().
		Background(theme.Surface).
		Width(s.width)

	modeSection := modeStyle.Render(strings.ToUpper(s.mode))
	available := s.width - lipgloss.Width(modeSection)
	if available < 1 {
		return barStyle.Render(cutToWidth(modeSection, s.width))
	}

	hints := buildHints(s.hints, hintKeyStyle, hintTextStyle, separatorStyle, available)
	hintsWidth := lipgloss.Width(hints)
	gap := available - hintsWidth
	if gap < 0 {
		gap = 0
	}

	line := lipgloss.JoinHorizontal(lipgloss.Center,
		modeSection,
		strings.Repeat(" ", gap),
		hints,
	)

	return barStyle.Render(cutToWidth(line, s.width))
}

func buildHints(hints []KeyHint, keyStyle, textStyle, sepStyle lipgloss.Style, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	sep := sepStyle.Render(" | ")
	parts := make([]string, 0, len(hints))
	for _, h := range hints {
		parts = append(parts, keyStyle.Render(h.Key)+" "+textStyle.Render(h.Description))
	}

	var chosen []string
	for _, p := range parts {
		candidate := p
		if len(chosen) > 0 {
			candidate = strings.Join(append(chosen, p), sep)
		}
		if lipgloss.Width(candidate) > maxWidth {
			break
		}
		chosen = append(chosen, p)
	}

	if len(chosen) == 0 {
		return ""
	}

	return strings.Join(chosen, sep)
}

func cutToWidth(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(value) <= width {
		return value
	}
	return lipgloss.NewStyle().Width(width).MaxWidth(width).Inline(true).Render(value)
}
