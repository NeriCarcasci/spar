package dashboard

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/challenge"
	"github.com/spar-cli/spar/internal/profile"
	"github.com/spar-cli/spar/internal/ui/components"
	"github.com/spar-cli/spar/internal/ui/theme"
)

type NavigateBrowserMsg struct{}

type NavigateRandomMsg struct{}

type Model struct {
	width     int
	height    int
	profile   *profile.Profile
	index     *challenge.Index
	statusBar components.StatusBar
}

func New(p *profile.Profile, idx *challenge.Index) Model {
	return Model{
		profile: p,
		index:   idx,
		statusBar: components.NewStatusBar().
			WithMode("DASHBOARD").
			WithHints([]components.KeyHint{
				{Key: "b", Description: "browse"},
				{Key: "r", Description: "random"},
				{Key: "q", Description: "quit"},
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
		case "b":
			return m, navigateBrowser
		case "r":
			return m, navigateRandom
		}
	}
	return m, nil
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.TextPrimary).
		Bold(true).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.TextMid)

	totalAvailable := len(m.index.Challenges)
	totalSolved := m.profile.TotalSolved()
	streak := m.profile.CurrentStreak()

	streakStyle := lipgloss.NewStyle().Bold(true)
	if streak > 0 {
		streakStyle = streakStyle.Foreground(theme.Red)
	} else {
		streakStyle = streakStyle.Foreground(theme.TextDim)
	}

	progressStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary)

	title := titleStyle.Render("spar")

	stats := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("Streak  ")+streakStyle.Render(formatStreak(streak)),
		labelStyle.Render("Solved  ")+progressStyle.Render(formatProgress(totalSolved, totalAvailable)),
		"",
		labelStyle.Render("press b to browse challenges or r for a random challenge"),
	)

	cardStyle := lipgloss.NewStyle().
		Background(theme.Surface).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 2)

	card := cardStyle.Render(stats)
	content := lipgloss.JoinVertical(lipgloss.Center, title, "", card)

	body := lipgloss.Place(
		m.width, m.height-1,
		lipgloss.Center, lipgloss.Center,
		content,
	)

	statusBar := m.statusBar.WithWidth(m.width).View()

	return lipgloss.JoinVertical(lipgloss.Left, body, statusBar)
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

func navigateBrowser() tea.Msg {
	return NavigateBrowserMsg{}
}

func navigateRandom() tea.Msg {
	return NavigateRandomMsg{}
}

func formatStreak(days int) string {
	if days == 1 {
		return "1 day"
	}
	return strconv.Itoa(days) + " days"
}

func formatProgress(solved, total int) string {
	return strconv.Itoa(solved) + " / " + strconv.Itoa(total)
}
