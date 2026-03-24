package profile

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	profiledata "github.com/NeriCarcasci/spar/internal/profile"
	"github.com/NeriCarcasci/spar/internal/ui/components"
	"github.com/NeriCarcasci/spar/internal/ui/theme"
)

type NavigateDashboardMsg struct{}

type Model struct {
	width     int
	height    int
	profile   *profiledata.Profile
	statusBar components.StatusBar
}

func New(p *profiledata.Profile) Model {
	return Model{
		profile: p,
		statusBar: components.NewStatusBar().
			WithMode("PROFILE").
			WithHints([]components.KeyHint{
				{Key: "esc", Description: "back"},
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
		case "esc":
			return m, navigateDashboard
		}
	}
	return m, nil
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.TextPrimary).
		Bold(true)

	statStyle := lipgloss.NewStyle().
		Foreground(theme.TextMid)

	header := titleStyle.Render("Profile")
	stats := statStyle.Render("full stats coming in milestone 7")

	cardStyle := lipgloss.NewStyle().
		Background(theme.Surface).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 2)

	card := cardStyle.Render(lipgloss.JoinVertical(lipgloss.Left, header, "", stats))

	placed := lipgloss.Place(m.width, m.height-1, lipgloss.Left, lipgloss.Top,
		lipgloss.NewStyle().Padding(1, 2).Render(card))

	statusBar := m.statusBar.WithWidth(m.width).View()
	return lipgloss.JoinVertical(lipgloss.Left, placed, statusBar)
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

func navigateDashboard() tea.Msg {
	return NavigateDashboardMsg{}
}
