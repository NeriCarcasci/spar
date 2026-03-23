package session

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/challenge"
	"github.com/spar-cli/spar/internal/ui/components"
	"github.com/spar-cli/spar/internal/ui/theme"
)

type NavigateBrowserMsg struct{}

type Model struct {
	width     int
	height    int
	challenge *challenge.Challenge
	language  string
	statusBar components.StatusBar
}

func New(ch *challenge.Challenge, language string) Model {
	return Model{
		challenge: ch,
		language:  language,
		statusBar: components.NewStatusBar().
			WithMode("SESSION").
			WithHints([]components.KeyHint{
				{Key: "esc", Description: "back"},
				{Key: "ctrl+r", Description: "run tests"},
				{Key: "ctrl+t", Description: "toggle AI"},
			}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, navigateBrowser
		}
	}
	return m, nil
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.TextPrimary).
		Bold(true)

	diffBadge := theme.DifficultyStyle(string(m.challenge.Difficulty))
	langBadge := theme.LanguageBadge()

	header := titleStyle.Render(m.challenge.Title) + " " +
		diffBadge.Render(string(m.challenge.Difficulty)) + " " +
		langBadge.Render(m.language)

	placeholder := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		Render("session view — editor and test runner coming in milestone 3+4")

	body := lipgloss.JoinVertical(lipgloss.Left, header, "", placeholder)

	placed := lipgloss.Place(m.width, m.height-1, lipgloss.Left, lipgloss.Top,
		lipgloss.NewStyle().Padding(1, 2).Render(body))

	statusBar := m.statusBar.WithWidth(m.width).View()
	return lipgloss.JoinVertical(lipgloss.Left, placed, statusBar)
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

func navigateBrowser() tea.Msg {
	return NavigateBrowserMsg{}
}
