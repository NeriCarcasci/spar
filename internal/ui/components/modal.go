package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/ui/theme"
)

type Modal struct {
	title   string
	body    string
	width   int
	height  int
	visible bool
}

func NewModal() Modal {
	return Modal{}
}

func (m Modal) WithTitle(title string) Modal {
	m.title = title
	return m
}

func (m Modal) WithBody(body string) Modal {
	m.body = body
	return m
}

func (m Modal) WithSize(width, height int) Modal {
	m.width = width
	m.height = height
	return m
}

func (m Modal) WithVisible(visible bool) Modal {
	m.visible = visible
	return m
}

func (m Modal) Visible() bool {
	return m.visible
}

func (m Modal) View() string {
	if !m.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Red).
		Padding(0, 1).
		MarginBottom(1)

	bodyStyle := lipgloss.NewStyle().
		Foreground(theme.TextPrimary).
		Padding(0, 1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Background(theme.Surface).
		Padding(1, 2)

	title := titleStyle.Render(m.title)
	body := bodyStyle.Render(m.body)
	content := lipgloss.JoinVertical(lipgloss.Left, title, body)

	box := boxStyle.
		Width(m.width - 4).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}
