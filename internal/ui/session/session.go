package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/challenge"
	"github.com/spar-cli/spar/internal/runner"
	"github.com/spar-cli/spar/internal/ui/components"
	"github.com/spar-cli/spar/internal/ui/theme"
)

const (
	tabChallenge   = 0
	tabCode        = 1
	tabInterviewer = 2
)

type NavigateBrowserMsg struct{}
type clearSavedMsg struct{}

type SavedMsg struct {
	Language string
	Err      error
}

type runDoneMsg struct {
	results    []runner.TestResult
	compileErr string
	runtimeErr string
	err        error
	duration   time.Duration
}

type chatRole int

const (
	roleInterviewer chatRole = iota
	roleUser
	roleSystem
)

type chatMessage struct {
	Role    chatRole
	Content string
}

type Model struct {
	width     int
	height    int
	challenge *challenge.Challenge
	language  string
	languages []string
	langIdx   int
	statusBar components.StatusBar
	activeTab int
	dirty     bool
	picking   bool
	running   bool
	saved     bool

	editor textarea.Model
	viewer viewport.Model

	chatInput    textarea.Model
	chatMessages []chatMessage
	chatView     viewport.Model

	results    []runner.TestResult
	compileErr string
	runtimeErr string
	lastRun    time.Duration
}

func New(ch *challenge.Challenge, language string) Model {
	editor := newEditor()
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle()

	chatIn := textarea.New()
	chatIn.Placeholder = "Ask the interviewer a question..."
	chatIn.ShowLineNumbers = false
	chatIn.CharLimit = 500
	chatIn.SetWidth(80)
	chatIn.SetHeight(3)
	chatIn.FocusedStyle.Base = lipgloss.NewStyle()
	chatIn.FocusedStyle.Text = lipgloss.NewStyle().Foreground(theme.TextPrimary)
	chatIn.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(theme.TextFaint)
	chatIn.FocusedStyle.CursorLine = lipgloss.NewStyle()
	chatIn.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(theme.Red)
	chatIn.BlurredStyle = chatIn.FocusedStyle
	chatIn.Prompt = "> "

	chatVP := viewport.New(80, 15)
	chatVP.Style = lipgloss.NewStyle()

	langs := ch.Languages
	idx := 0
	for i, l := range langs {
		if l == language {
			idx = i
			break
		}
	}

	intro := chatMessage{
		Role: roleInterviewer,
		Content: fmt.Sprintf(
			"Welcome! I'll be your interviewer today. You're working on \"%s\" (%s).\n\n"+
				"Take your time to read the problem, then switch to the Code tab to write your solution. "+
				"When you're ready, type /run to test or /submit to submit your solution for review.\n\n"+
				"Feel free to ask me questions — I can give hints, clarify constraints, or discuss your approach.",
			ch.Title, string(ch.Difficulty)),
	}

	return Model{
		challenge:    ch,
		language:     language,
		languages:    langs,
		langIdx:      idx,
		activeTab:    tabChallenge,
		editor:       editor,
		viewer:       vp,
		chatInput:    chatIn,
		chatMessages: []chatMessage{intro},
		chatView:     chatVP,
		statusBar: components.NewStatusBar().
			WithMode("SESSION").
			WithHints(sessionHints(tabChallenge)),
	}
}

func newEditor() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Write your solution here..."
	ta.ShowLineNumbers = true
	ta.CharLimit = 0
	ta.SetWidth(80)
	ta.SetHeight(20)
	ta.Focus()

	ta.Cursor.Style = lipgloss.NewStyle().Foreground(theme.Red)
	ta.FocusedStyle.Base = lipgloss.NewStyle()
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(theme.Surface)
	ta.FocusedStyle.LineNumber = lipgloss.NewStyle().Foreground(theme.TextFaint)
	ta.FocusedStyle.CursorLineNumber = lipgloss.NewStyle().Foreground(theme.TextDim)
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(theme.TextFaint)
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(theme.TextPrimary)
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(theme.TextDim)
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle().Foreground(theme.TextFaint)
	ta.BlurredStyle = ta.FocusedStyle

	return ta
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SavedMsg:
		if msg.Err == nil {
			m.dirty = false
			m.saved = true
			return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearSavedMsg{} })
		}
		return m, nil

	case clearSavedMsg:
		m.saved = false
		return m, nil

	case runDoneMsg:
		m.running = false
		m.results = msg.results
		m.compileErr = msg.compileErr
		m.runtimeErr = msg.runtimeErr
		m.lastRun = msg.duration

		summary := m.buildRunSummary(msg)
		m.chatMessages = append(m.chatMessages, chatMessage{Role: roleSystem, Content: summary})
		m.refreshChat()
		m.statusBar = m.statusBar.WithHints(sessionHints(m.activeTab))
		return m, nil

	case tea.KeyMsg:
		if m.picking {
			return m.updatePicker(msg)
		}

		switch msg.String() {
		case "esc":
			if m.dirty {
				return m, tea.Batch(m.saveCode(), navigateBrowser)
			}
			return m, navigateBrowser

		case "ctrl+l":
			m.picking = true
			m.statusBar = m.statusBar.WithHints([]components.KeyHint{
				{Key: "←/→", Description: "choose"},
				{Key: "enter", Description: "confirm"},
				{Key: "esc", Description: "cancel"},
			})
			return m, nil

		case "ctrl+s":
			return m, m.saveCode()

		case "ctrl+r":
			if m.running {
				return m, nil
			}
			m.running = true
			prev := m.activeTab
			m.activeTab = tabInterviewer
			m.updateFocus(prev)
			m.chatMessages = append(m.chatMessages, chatMessage{Role: roleSystem, Content: "Running tests..."})
			m.refreshChat()
			m.statusBar = m.statusBar.WithHints(sessionHints(m.activeTab))
			if m.dirty {
				return m, tea.Batch(m.saveCode(), m.runTests())
			}
			return m, m.runTests()

		case "f1":
			return m.switchTab(tabChallenge)
		case "f2":
			return m.switchTab(tabCode)
		case "f3":
			return m.switchTab(tabInterviewer)
		}
	}

	var cmd tea.Cmd
	switch m.activeTab {
	case tabCode:
		old := m.editor.Value()
		m.editor, cmd = m.editor.Update(msg)
		if m.editor.Value() != old {
			m.dirty = true
		}
	case tabInterviewer:
		cmd = m.updateChat(msg)
	default:
		m.viewer, cmd = m.viewer.Update(msg)
	}
	return m, cmd
}

func (m *Model) updateChat(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
		text := strings.TrimSpace(m.chatInput.Value())
		if text == "" {
			return nil
		}
		m.chatInput.Reset()

		if text == "/run" {
			m.chatMessages = append(m.chatMessages, chatMessage{Role: roleUser, Content: text})
			m.running = true
			m.chatMessages = append(m.chatMessages, chatMessage{Role: roleSystem, Content: "Running tests..."})
			m.refreshChat()
			if m.dirty {
				return tea.Batch(m.saveCode(), m.runTests())
			}
			return m.runTests()
		}

		if text == "/submit" {
			m.chatMessages = append(m.chatMessages, chatMessage{Role: roleUser, Content: text})
			m.chatMessages = append(m.chatMessages, chatMessage{Role: roleInterviewer,
				Content: "Let me review your solution...\n\n" +
					"[AI review will be available when an AI provider is configured in Settings. " +
					"For now, use /run to test against the test cases.]"})
			m.refreshChat()
			return nil
		}

		m.chatMessages = append(m.chatMessages, chatMessage{Role: roleUser, Content: text})
		m.chatMessages = append(m.chatMessages, chatMessage{Role: roleInterviewer,
			Content: "[AI responses will be available when an AI provider is configured in Settings. " +
				"You can still use /run to test your code and /submit to submit.]"})
		m.refreshChat()
		return nil
	}

	var cmd tea.Cmd
	m.chatInput, cmd = m.chatInput.Update(msg)
	return cmd
}

func (m Model) switchTab(tab int) (Model, tea.Cmd) {
	if tab == m.activeTab {
		return m, nil
	}
	prev := m.activeTab
	m.activeTab = tab
	m.updateFocus(prev)
	m.statusBar = m.statusBar.WithHints(sessionHints(m.activeTab))
	if tab == tabCode || tab == tabInterviewer {
		return m, textarea.Blink
	}
	return m, nil
}

func (m *Model) updateFocus(prevTab int) {
	switch prevTab {
	case tabCode:
		m.editor.Blur()
	case tabInterviewer:
		m.chatInput.Blur()
	}
	switch m.activeTab {
	case tabCode:
		m.editor.Focus()
	case tabInterviewer:
		m.chatInput.Focus()
	}
}

func (m Model) updatePicker(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		m.langIdx = (m.langIdx + len(m.languages) - 1) % len(m.languages)
		return m, nil
	case "right", "l":
		m.langIdx = (m.langIdx + 1) % len(m.languages)
		return m, nil
	case "enter":
		newLang := m.languages[m.langIdx]
		m.picking = false
		m.statusBar = m.statusBar.WithHints(sessionHints(m.activeTab))
		if newLang == m.language {
			return m, nil
		}
		var cmds []tea.Cmd
		if m.dirty {
			cmds = append(cmds, m.saveCode())
		}
		m.language = newLang
		code, _ := challenge.LoadSetupCode(m.challenge.Path, newLang)
		m.editor.SetValue(code)
		m.editor.CursorStart()
		m.dirty = false
		return m, tea.Batch(cmds...)
	case "esc":
		m.picking = false
		for i, l := range m.languages {
			if l == m.language {
				m.langIdx = i
				break
			}
		}
		m.statusBar = m.statusBar.WithHints(sessionHints(m.activeTab))
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	header := m.renderHeader()
	headerH := lipgloss.Height(header)
	statusBar := m.statusBar.WithWidth(m.width).View()
	statusH := lipgloss.Height(statusBar)
	bodyH := max(1, m.height-headerH-statusH)

	var body string
	switch m.activeTab {
	case tabChallenge:
		body = m.renderChallenge(m.width, bodyH)
	case tabCode:
		body = m.renderCode(m.width, bodyH)
	case tabInterviewer:
		body = m.renderInterviewer(m.width, bodyH)
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, body, statusBar)
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height

	if m.challenge == nil {
		return m
	}

	contentW := max(1, width-4)
	contentH := max(1, height-6)

	m.editor.SetWidth(contentW)
	m.editor.SetHeight(contentH)
	m.viewer.Width = contentW
	m.viewer.Height = contentH
	m.viewer.SetContent(m.challengeContent(contentW))

	m.chatInput.SetWidth(contentW)
	m.chatView.Width = contentW
	m.chatView.Height = max(1, contentH-5)
	m.refreshChat()

	return m
}

func (m Model) renderHeader() string {
	title := lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render(m.challenge.Title)
	diff := theme.DifficultyStyle(string(m.challenge.Difficulty)).Render(string(m.challenge.Difficulty))

	var lang string
	if m.picking {
		langParts := make([]string, len(m.languages))
		for i, l := range m.languages {
			if i == m.langIdx {
				langParts[i] = lipgloss.NewStyle().
					Foreground(theme.TextPrimary).Background(theme.Surface2).Bold(true).Padding(0, 1).Render(l)
			} else {
				langParts[i] = lipgloss.NewStyle().Foreground(theme.TextFaint).Padding(0, 1).Render(l)
			}
		}
		lang = strings.Join(langParts, " ")
	} else {
		lang = theme.LanguageBadge().Render(m.language)
	}

	dirtyMark := ""
	if m.saved {
		dirtyMark = lipgloss.NewStyle().Foreground(theme.Green).Render(" saved")
	} else if m.dirty {
		dirtyMark = lipgloss.NewStyle().Foreground(theme.Amber).Render(" *")
	}

	headerLine := title + " " + diff + " " + lang + dirtyMark

	tabNames := []string{"Challenge", "Code", "Interviewer"}
	tabParts := make([]string, len(tabNames))
	for i, t := range tabNames {
		if i == m.activeTab {
			tabParts[i] = lipgloss.NewStyle().
				Foreground(theme.Background).
				Background(theme.Red).
				Bold(true).
				Padding(0, 1).
				Render(t)
		} else {
			tabParts[i] = lipgloss.NewStyle().Foreground(theme.TextDim).Padding(0, 1).Render(t)
		}
	}
	sep := lipgloss.NewStyle().Foreground(theme.Border).Render("│")
	tabLine := strings.Join(tabParts, sep)
	divider := lipgloss.NewStyle().Foreground(theme.Border).Render(strings.Repeat("─", max(1, m.width-4)))

	return lipgloss.NewStyle().Width(m.width).Padding(0, 2).Render(
		lipgloss.JoinVertical(lipgloss.Left, headerLine, "", tabLine, divider),
	)
}

func (m Model) renderChallenge(width, height int) string {
	return lipgloss.NewStyle().Width(width).Height(height).Padding(0, 2).Render(m.viewer.View())
}

func (m Model) renderCode(width, height int) string {
	return lipgloss.NewStyle().Width(width).Height(height).Padding(0, 2).Render(m.editor.View())
}

func (m Model) renderInterviewer(width, height int) string {
	contentW := max(1, width-4)
	inputH := 3
	chatH := max(1, height-inputH-1)

	m.chatView.Width = contentW
	m.chatView.Height = chatH

	inputBorder := lipgloss.NewStyle().Foreground(theme.Border).Render(strings.Repeat("─", contentW))

	return lipgloss.NewStyle().Width(width).Height(height).Padding(0, 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			m.chatView.View(),
			inputBorder,
			m.chatInput.View(),
		),
	)
}

func (m *Model) refreshChat() {
	contentW := max(1, m.width-6)
	var lines []string
	for _, msg := range m.chatMessages {
		var styled string
		switch msg.Role {
		case roleInterviewer:
			label := lipgloss.NewStyle().Foreground(theme.Red).Bold(true).Render("Interviewer")
			body := lipgloss.NewStyle().Foreground(theme.TextMid).Width(contentW).Render(msg.Content)
			styled = label + "\n" + body
		case roleUser:
			label := lipgloss.NewStyle().Foreground(theme.Green).Bold(true).Render("You")
			body := lipgloss.NewStyle().Foreground(theme.TextPrimary).Width(contentW).Render(msg.Content)
			styled = label + "\n" + body
		case roleSystem:
			boxW := max(1, contentW-4)
			body := lipgloss.NewStyle().Foreground(theme.Amber).Width(boxW).Render(msg.Content)
			styled = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.AmberDim).
				Padding(0, 1).
				Width(boxW).
				Render(body)
		}
		lines = append(lines, styled)
	}
	m.chatView.SetContent(strings.Join(lines, "\n\n"))
	m.chatView.GotoBottom()
}

func (m Model) challengeContent(width int) string {
	ch := m.challenge
	dimStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
	headStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true)
	bodyStyle := lipgloss.NewStyle().Foreground(theme.TextMid).Width(max(1, width-2))

	var sections []string

	sections = append(sections,
		headStyle.Render("Description"),
		bodyStyle.Render(strings.TrimSpace(ch.Description)),
		"",
	)

	if len(ch.Constraints) > 0 {
		constraints := make([]string, len(ch.Constraints))
		for i, c := range ch.Constraints {
			constraints[i] = dimStyle.Render("  • ") + bodyStyle.Render(c)
		}
		sections = append(sections,
			headStyle.Render("Constraints"),
			strings.Join(constraints, "\n"),
			"",
		)
	}

	if len(ch.Examples) > 0 {
		sections = append(sections, headStyle.Render("Examples"))
		for i, ex := range ch.Examples {
			sections = append(sections,
				dimStyle.Render(fmt.Sprintf("  Example %d:", i+1)),
				dimStyle.Render("    Input:  ")+bodyStyle.Render(ex.Input),
				dimStyle.Render("    Output: ")+bodyStyle.Render(ex.Output),
			)
			if ex.Explanation != "" {
				sections = append(sections, dimStyle.Render("    → ")+bodyStyle.Render(ex.Explanation))
			}
		}
		sections = append(sections, "")
	}

	if len(ch.Hints) > 0 {
		sections = append(sections, headStyle.Render("Hints"))
		for i, h := range ch.Hints {
			sections = append(sections,
				dimStyle.Render(fmt.Sprintf("  %d. ", i+1))+bodyStyle.Render(h),
			)
		}
	}

	return strings.Join(sections, "\n")
}

func (m Model) runTests() tea.Cmd {
	chDir := m.challenge.Path
	lang := m.language
	code := m.editor.Value()
	return func() tea.Msg {
		start := time.Now()
		results, compileErr, runtimeErr, err := runner.Run(chDir, lang, code)
		return runDoneMsg{
			results:    results,
			compileErr: compileErr,
			runtimeErr: runtimeErr,
			err:        err,
			duration:   time.Since(start),
		}
	}
}

func (m Model) buildRunSummary(msg runDoneMsg) string {
	if msg.err != nil {
		errStr := msg.err.Error()
		if strings.Contains(errStr, "cannot find") || strings.Contains(errStr, "no such file") {
			return "Builder not found for this challenge.\n" +
				"The builder/ folder needs builder files to run tests.\n" +
				"See the project README for how to set up builders."
		}
		return "Error: " + errStr
	}
	if msg.compileErr != "" {
		return "Compile error:\n" + msg.compileErr
	}
	if msg.runtimeErr != "" {
		return "Runtime error:\n" + msg.runtimeErr
	}
	if len(msg.results) == 0 {
		return "No test results (builder may not be configured for this challenge yet)"
	}

	passed, total := 0, len(msg.results)
	for _, r := range msg.results {
		if r.Passed {
			passed++
		}
	}

	var b strings.Builder
	for _, r := range msg.results {
		if r.Passed {
			b.WriteString(fmt.Sprintf("  ✓ Test %d passed\n", r.Index))
		} else {
			b.WriteString(fmt.Sprintf("  ✗ Test %d failed — got %s, expected %s\n", r.Index, r.Got, r.Expected))
		}
	}

	status := fmt.Sprintf("Results: %d/%d passed (%s)", passed, total, msg.duration.Round(time.Millisecond))
	if passed == total {
		status = fmt.Sprintf("All %d tests passed! (%s) 🎉", total, msg.duration.Round(time.Millisecond))
	}

	return status + "\n" + b.String()
}

func (m Model) saveCode() tea.Cmd {
	lang := m.language
	code := m.editor.Value()
	chPath := m.challenge.Path
	return func() tea.Msg {
		filename := setupFilename(lang)
		path := filepath.Join(chPath, "setup", filename)
		err := os.WriteFile(path, []byte(code), 0644)
		return SavedMsg{Language: lang, Err: err}
	}
}

func (m Model) WithCode(code string) Model {
	m.editor.SetValue(code)
	m.editor.CursorStart()
	return m
}

func sessionHints(tab int) []components.KeyHint {
	base := []components.KeyHint{
		{Key: "f1", Description: "challenge"},
		{Key: "f2", Description: "code"},
		{Key: "f3", Description: "interviewer"},
		{Key: "ctrl+r", Description: "run"},
		{Key: "ctrl+s", Description: "save"},
		{Key: "esc", Description: "back"},
	}
	if tab == tabInterviewer {
		return []components.KeyHint{
			{Key: "f1/f2/f3", Description: "tabs"},
			{Key: "/run", Description: "test"},
			{Key: "/submit", Description: "review"},
			{Key: "esc", Description: "back"},
		}
	}
	return base
}

func setupFilename(language string) string {
	switch language {
	case "python":
		return "python.py"
	case "go":
		return "go.go"
	case "javascript":
		return "javascript.js"
	case "cpp":
		return "cpp.cpp"
	case "rust":
		return "rust.rs"
	default:
		return language + ".txt"
	}
}

func navigateBrowser() tea.Msg {
	return NavigateBrowserMsg{}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
