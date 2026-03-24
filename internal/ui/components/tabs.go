package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/NeriCarcasci/spar/internal/ui/theme"
)

type Tabs struct {
	items       []string
	activeIndex int
	width       int
}

func NewTabs(items []string) Tabs {
	return Tabs{items: items}
}

func (t Tabs) WithActive(index int) Tabs {
	if index >= 0 && index < len(t.items) {
		t.activeIndex = index
	}
	return t
}

func (t Tabs) WithWidth(width int) Tabs {
	t.width = width
	return t
}

func (t Tabs) ActiveIndex() int {
	return t.activeIndex
}

func (t Tabs) View() string {
	activeStyle := lipgloss.NewStyle().
		Foreground(theme.Red).
		Bold(true).
		Padding(0, 2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(theme.Red)

	normalStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		Padding(0, 2)

	sepStyle := lipgloss.NewStyle().
		Foreground(theme.Border)

	var rendered []string
	for i, item := range t.items {
		if i == t.activeIndex {
			rendered = append(rendered, activeStyle.Render(item))
		} else {
			rendered = append(rendered, normalStyle.Render(item))
		}
	}

	row := strings.Join(rendered, sepStyle.Render(" │ "))

	return lipgloss.NewStyle().Width(t.width).Render(row)
}
