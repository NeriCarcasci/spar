package splash

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/ui/theme"
)

const minimumDisplayDuration = 1500 * time.Millisecond

const logo = `
  ____  ____   ___   ____
 / ___||  _ \ / _ \ |  _ \
 \___ \| |_) | | | || |_) |
  ___) |  __/| |_| ||  _ <
 |____/|_|    \___/ |_| \_\
`

const tagline = "code under pressure"

type LoadingCompleteMsg struct{}

type MinimumTimeElapsedMsg struct{}

type Model struct {
	width              int
	height             int
	loadingDone        bool
	minimumTimeElapsed bool
	frame              int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		startMinimumTimer(),
		animateTick(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case MinimumTimeElapsedMsg:
		m.minimumTimeElapsed = true
	case LoadingCompleteMsg:
		m.loadingDone = true
	case tickMsg:
		m.frame++
		if !m.ReadyToTransition() {
			return m, animateTick()
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	logoStyle := lipgloss.NewStyle().
		Foreground(theme.Red).
		Bold(true)

	taglineStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		Italic(true)

	spinnerStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logoStyle.Render(logo),
		"",
		taglineStyle.Render(tagline),
		"",
		spinnerStyle.Render(spinnerFrame(m.frame)),
	)

	placed := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(theme.Background).
		Render(placed)
}

func (m Model) ReadyToTransition() bool {
	return m.loadingDone && m.minimumTimeElapsed
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

type tickMsg struct{}

func animateTick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func startMinimumTimer() tea.Cmd {
	return tea.Tick(minimumDisplayDuration, func(time.Time) tea.Msg {
		return MinimumTimeElapsedMsg{}
	})
}

func spinnerFrame(frame int) string {
	frames := []string{"-", "\\", "|", "/"}
	return frames[frame%len(frames)] + " loading..."
}
