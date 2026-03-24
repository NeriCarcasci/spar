package browser

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/challenge"
	"github.com/spar-cli/spar/internal/profile"
	"github.com/spar-cli/spar/internal/ui/components"
	"github.com/spar-cli/spar/internal/ui/theme"
)

type SelectChallengeMsg struct {
	Entry challenge.IndexEntry
}

type NavigateDashboardMsg struct{}

type Model struct {
	width     int
	height    int
	index     *challenge.Index
	profile   *profile.Profile
	filtered  []challenge.IndexEntry
	filters   Filters
	cursor    int
	offset    int
	statusBar components.StatusBar
}

func New(idx *challenge.Index, p *profile.Profile) Model {
	m := Model{
		index:   idx,
		profile: p,
		statusBar: components.NewStatusBar().
			WithMode("BROWSE").
			WithHints([]components.KeyHint{
				{Key: "↑↓", Description: "navigate"},
				{Key: "enter", Description: "select"},
				{Key: "esc", Description: "back"},
				{Key: "q", Description: "quit"},
			}),
	}
	m.applyFilters()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveCursorUp()
		case "down", "j":
			m.moveCursorDown()
		case "enter":
			if len(m.filtered) > 0 {
				return m, selectChallenge(m.filtered[m.cursor])
			}
		case "esc":
			return m, navigateDashboard
		}
	}
	return m, nil
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Red).
		Bold(true)

	header := titleStyle.Render("Challenges")

	if len(m.filtered) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
		body := lipgloss.JoinVertical(lipgloss.Left, header, "", emptyStyle.Render("No challenges found"))
		return m.wrapWithStatusBar(body)
	}

	visibleHeight := m.height - 5
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	var rows []string
	end := m.offset + visibleHeight
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := m.offset; i < end; i++ {
		rows = append(rows, m.renderRow(i))
	}

	list := strings.Join(rows, "\n")
	body := lipgloss.JoinVertical(lipgloss.Left, header, "", list)
	return m.wrapWithStatusBar(body)
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

func (m *Model) applyFilters() {
	isSolved := func(id string) bool {
		return m.profile.IsSolved(id)
	}
	m.filtered = m.filters.Apply(m.index.Challenges, isSolved)
	m.cursor = 0
	m.offset = 0
}

func (m *Model) moveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
	}
}

func (m *Model) moveCursorDown() {
	if m.cursor < len(m.filtered)-1 {
		m.cursor++
		visibleHeight := m.height - 5
		if visibleHeight < 1 {
			visibleHeight = 1
		}
		if m.cursor >= m.offset+visibleHeight {
			m.offset = m.cursor - visibleHeight + 1
		}
	}
}

func (m Model) renderRow(index int) string {
	entry := m.filtered[index]
	isCursor := index == m.cursor
	solved := m.profile.IsSolved(entry.ID)

	diffBadge := theme.DifficultyStyle(string(entry.Difficulty)).Render(padRight(string(entry.Difficulty), 6))
	categoryBadge := theme.CategoryBadge().Render(entry.Category)

	var titleStyle lipgloss.Style
	var prefix string
	if isCursor {
		prefix = "▸ "
		titleStyle = lipgloss.NewStyle().
			Foreground(theme.TextPrimary).
			Background(theme.Surface2).
			Bold(true)
	} else {
		prefix = "  "
		titleStyle = lipgloss.NewStyle().
			Foreground(theme.TextMid)
	}

	solvedMark := lipgloss.NewStyle().Foreground(theme.TextDim).Render("○")
	if solved {
		solvedMark = lipgloss.NewStyle().Foreground(theme.Green).Render("✓")
	}

	titleWidth := m.width - 24
	if titleWidth < 12 {
		titleWidth = 12
	}
	if titleWidth > 40 {
		titleWidth = 40
	}
	return prefix + solvedMark + " " +
		diffBadge + " " +
		titleStyle.Render(padRight(entry.Title, titleWidth)) + " " +
		categoryBadge
}

func (m Model) wrapWithStatusBar(body string) string {
	placed := lipgloss.Place(m.width, m.height-1, lipgloss.Left, lipgloss.Top,
		lipgloss.NewStyle().Padding(1, 2).Render(body))
	statusBar := m.statusBar.WithWidth(m.width).View()
	return lipgloss.JoinVertical(lipgloss.Left, placed, statusBar)
}

func selectChallenge(entry challenge.IndexEntry) tea.Cmd {
	return func() tea.Msg {
		return SelectChallengeMsg{Entry: entry}
	}
}

func navigateDashboard() tea.Msg {
	return NavigateDashboardMsg{}
}

func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		runes := []rune(s)
		for len(runes) > 0 && lipgloss.Width(string(runes)) > width {
			runes = runes[:len(runes)-1]
		}
		return string(runes)
	}
	return s + strings.Repeat(" ", width-w)
}
