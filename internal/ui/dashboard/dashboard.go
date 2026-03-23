package dashboard

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/challenge"
	"github.com/spar-cli/spar/internal/config"
	"github.com/spar-cli/spar/internal/profile"
	"github.com/spar-cli/spar/internal/ui/components"
	"github.com/spar-cli/spar/internal/ui/theme"
)

const (
	viewDashboard = "dashboard"
	viewBrowse    = "browse"
	viewProfile   = "profile"
	viewRandom    = "random"
	viewSettings  = "settings"
)

const sidebarLogo = "SPAR\n----"

type artTickMsg struct{}

type SelectChallengeMsg struct{ Entry challenge.IndexEntry }

type ConfigChangedMsg struct{ Config config.Config }

type collection struct {
	Name   string
	Solved int
	Total  int
}

type setting struct {
	Key     string
	Value   string
	Options []string
	Index   int
}

type Model struct {
	width, height int
	profile       *profile.Profile
	index         *challenge.Index
	cfg           config.Config
	sidebar       components.Sidebar
	statusBar     components.StatusBar
	view          string
	artFrame      int

	collections   []collection
	browseCursor  int
	browseOpen    bool
	challengeCur  int
	challengeOff  int
	activeCollect string

	randomFilter int

	settings       []setting
	settingCur     int
	editingSetting bool
}

func New(p *profile.Profile, idx *challenge.Index, cfg config.Config) Model {
	items := []components.SidebarItem{
		{ID: viewDashboard, Label: "Dashboard", Key: "d"},
		{ID: viewBrowse, Label: "Browse", Key: "b"},
		{ID: viewProfile, Label: "Profile", Key: "p"},
		{ID: viewRandom, Label: "Random", Key: "r"},
		{ID: viewSettings, Label: "Settings", Key: "s", Secondary: true},
		{ID: components.SidebarQuitID, Label: "Quit", Key: "q", Secondary: true},
	}
	m := Model{
		profile:   p,
		index:     idx,
		cfg:       cfg,
		sidebar:   components.NewSidebar(items, viewDashboard),
		statusBar: components.NewStatusBar(),
		view:      viewDashboard,
		artFrame:  3,
	}
	m.refreshCollections()
	m.resetSettings()
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case artTickMsg:
		if m.artFrame < 3 {
			m.artFrame++
			return m, animateArtTick()
		}
		return m, nil
	case tea.KeyMsg:
		if shortcut(msg.String()) {
			return m.applySidebar(msg)
		}
		if handled, cmd := m.handleViewKey(msg); handled {
			return m, cmd
		}
		return m.applySidebar(msg)
	}
	return m.applySidebar(msg)
}

func (m Model) applySidebar(msg tea.Msg) (Model, tea.Cmd) {
	sb, ev := m.sidebar.Update(msg)
	m.sidebar = sb
	if ev.Quit {
		return m, tea.Quit
	}
	if ev.Selected.ID == "" {
		return m, nil
	}
	m.view = ev.Selected.ID
	if ev.SelectedChanged {
		m.browseOpen = false
		m.editingSetting = false
		m.artFrame = 0
		return m, animateArtTick()
	}
	return m, nil
}

func (m *Model) handleViewKey(msg tea.KeyMsg) (bool, tea.Cmd) {
	s := msg.String()
	switch m.view {
	case viewBrowse:
		if !m.browseOpen {
			switch s {
			case "left", "h", "up", "k":
				if m.browseCursor > 0 {
					m.browseCursor--
				}
				return true, nil
			case "right", "l", "down", "j":
				if m.browseCursor < len(m.collections)-1 {
					m.browseCursor++
				}
				return true, nil
			case "enter":
				if len(m.collections) == 0 {
					return true, nil
				}
				m.browseOpen = true
				m.activeCollect = m.collections[m.browseCursor].Name
				m.challengeCur, m.challengeOff = 0, 0
				return true, nil
			}
			return false, nil
		}
		entries := m.collectionEntries()
		switch s {
		case "esc":
			m.browseOpen = false
			return true, nil
		case "up", "k":
			if m.challengeCur > 0 {
				m.challengeCur--
			}
			if m.challengeCur < m.challengeOff {
				m.challengeOff = m.challengeCur
			}
			return true, nil
		case "down", "j":
			if m.challengeCur < len(entries)-1 {
				m.challengeCur++
			}
			vis := max(1, m.height-8)
			if m.challengeCur >= m.challengeOff+vis {
				m.challengeOff = m.challengeCur - vis + 1
			}
			return true, nil
		case "enter":
			if len(entries) == 0 {
				return true, nil
			}
			return true, startChallenge(entries[m.challengeCur])
		}
	case viewRandom:
		switch s {
		case "left", "h":
			m.randomFilter = (m.randomFilter + 3) % 4
			return true, nil
		case "right", "l":
			m.randomFilter = (m.randomFilter + 1) % 4
			return true, nil
		case "enter":
			c := m.randomEntries()
			if len(c) == 0 {
				return true, nil
			}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			return true, startChallenge(c[r.Intn(len(c))])
		}
	case viewSettings:
		if len(m.settings) == 0 {
			return false, nil
		}
		if !m.editingSetting {
			switch s {
			case "up", "k":
				if m.settingCur > 0 {
					m.settingCur--
				}
				return true, nil
			case "down", "j":
				if m.settingCur < len(m.settings)-1 {
					m.settingCur++
				}
				return true, nil
			case "enter":
				if len(m.settings[m.settingCur].Options) > 0 {
					m.editingSetting = true
					return true, nil
				}
			}
			return false, nil
		}
		cur := &m.settings[m.settingCur]
		switch s {
		case "esc":
			m.editingSetting = false
			m.resetSettings()
			return true, nil
		case "left", "h", "up", "k":
			if len(cur.Options) > 0 {
				cur.Index = (cur.Index + len(cur.Options) - 1) % len(cur.Options)
			}
			return true, nil
		case "right", "l", "down", "j":
			if len(cur.Options) > 0 {
				cur.Index = (cur.Index + 1) % len(cur.Options)
			}
			return true, nil
		case "enter":
			m.applySetting(cur)
			m.editingSetting = false
			m.resetSettings()
			return true, configChanged(m.cfg)
		}
	}
	return false, nil
}

func (m Model) SetSize(width, height int) Model { m.width, m.height = width, height; return m }
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	bodyH := max(1, m.height-1)
	sideW := m.sidebarWidth()
	contentW := max(1, m.width-sideW)

	side := m.sidebar.View(max(8, sideW-1), bodyH, sidebarLogo)
	content := m.renderContent(contentW, bodyH)
	status := m.statusBar.WithWidth(m.width).WithMode(m.mode()).WithHints(m.hints()).View()

	full := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, side, content),
		status,
	)
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Background(theme.Background).Render(full)
}

func (m Model) renderContent(width, height int) string {
	innerW, innerH := max(1, width-4), max(1, height-2)
	main := trimHeight(m.renderView(innerW, innerH), innerH)
	free := innerH - lipgloss.Height(main)
	if free < 0 {
		free = 0
	}
	if free > 0 {
		main += "\n" + m.renderArt(innerW, free)
	}
	return lipgloss.NewStyle().Width(width).Height(height).Padding(1, 2).Render(main)
}

func (m Model) renderView(width, height int) string {
	switch m.view {
	case viewBrowse:
		return m.renderBrowse(width, height)
	case viewProfile:
		return m.renderProfile(width)
	case viewRandom:
		return m.renderRandom(width, height)
	case viewSettings:
		return m.renderSettings(width)
	default:
		return m.renderDashboard(width)
	}
}

func (m Model) renderDashboard(width int) string {
	streak := m.profile.CurrentStreak()
	streakColor := theme.TextDim
	if streak > 0 {
		streakColor = theme.Red
	}
	cards := lipgloss.JoinHorizontal(lipgloss.Top,
		statCard("Streak", formatStreak(streak), streakColor, max(16, (width-4)/3)), "  ",
		statCard("Solved", strconv.Itoa(m.profile.TotalSolved())+" / "+strconv.Itoa(len(m.index.Challenges)), theme.TextPrimary, max(16, (width-4)/3)), "  ",
		statCard("Avg solve", m.avgSolve(), theme.TextPrimary, max(16, (width-4)/3)),
	)
	rows := m.recentRows(width)
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render("Dashboard"), "",
		cards, "",
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render("Recent sessions"),
		rows, "",
		lipgloss.NewStyle().Foreground(theme.TextDim).Render("press enter on Browse to find challenges"),
	)
}

func (m Model) renderBrowse(width, height int) string {
	head := lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render("Collections")
	if len(m.collections) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, head, "", "No collections")
	}
	if !m.browseOpen {
		return lipgloss.JoinVertical(lipgloss.Left, head, "", m.collectionGrid(width), "", lipgloss.NewStyle().Foreground(theme.TextDim).Render("enter opens selected collection"))
	}
	entries := m.collectionEntries()
	if len(entries) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, head, "", "No challenges", "", lipgloss.NewStyle().Foreground(theme.TextDim).Render("esc back"))
	}
	vis := max(1, height-6)
	start := min(max(0, m.challengeOff), max(0, len(entries)-vis))
	end := min(len(entries), start+vis)
	rows := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		rows = append(rows, m.challengeRow(entries[i], width, i == m.challengeCur))
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render("Collection: "+m.activeCollect), "",
		strings.Join(rows, "\n"), "",
		lipgloss.NewStyle().Foreground(theme.TextDim).Render("esc back"),
	)
}

func (m Model) renderProfile(width int) string {
	name := strings.TrimSpace(m.profile.Username)
	if name == "" {
		name = "anonymous"
	}
	avatar := lipgloss.NewStyle().Background(theme.Red).Foreground(theme.Background).Bold(true).Padding(1, 2).Render(strings.ToUpper(string(name[0])))
	header := lipgloss.JoinHorizontal(lipgloss.Center, avatar, "  ", lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render(name),
		lipgloss.NewStyle().Foreground(theme.TextDim).Render(strconv.Itoa(m.profile.TotalSolved())+" challenges solved"),
	))
	e, md, h := m.diffCounts()
	langs := m.languagePills()
	cats := m.categoryBars(width)
	return lipgloss.JoinVertical(lipgloss.Left,
		header, "",
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render("Languages"),
		langs, "",
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render("Difficulty breakdown"),
		lipgloss.JoinHorizontal(lipgloss.Top,
			statCard("Easy", strconv.Itoa(e), theme.Green, max(12, (width-4)/3)), "  ",
			statCard("Medium", strconv.Itoa(md), theme.Amber, max(12, (width-4)/3)), "  ",
			statCard("Hard", strconv.Itoa(h), theme.Red, max(12, (width-4)/3)),
		), "",
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render("Category breakdown"),
		cats,
	)
}

func (m Model) renderRandom(width, height int) string {
	labels := []string{"all", "easy", "medium", "hard"}
	parts := make([]string, 0, 4)
	for i, l := range labels {
		st := lipgloss.NewStyle().Foreground(theme.TextDim).Padding(0, 1)
		if i == m.randomFilter {
			st = st.Foreground(theme.TextPrimary).Background(theme.Surface2)
		}
		parts = append(parts, st.Render(l))
	}
	block := lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render("Random"), "",
		lipgloss.NewStyle().Foreground(theme.TextMid).Render("Press enter to start a random challenge"), "",
		lipgloss.NewStyle().Foreground(theme.TextDim).Render("Difficulty"),
		strings.Join(parts, " "), "",
		lipgloss.NewStyle().Foreground(theme.TextDim).Render(strconv.Itoa(len(m.randomEntries()))+" matching challenges"),
	)
	return lipgloss.Place(width, max(1, height-1), lipgloss.Center, lipgloss.Center, block)
}

func (m Model) renderSettings(width int) string {
	rows := make([]string, 0, len(m.settings))
	for i, s := range m.settings {
		rows = append(rows, m.settingRow(s, width, i == m.settingCur))
	}
	foot := "j/k navigate | enter edit"
	if m.editingSetting {
		foot = "left/right choose | enter apply | esc cancel"
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render("Settings"), "",
		strings.Join(rows, "\n"), "",
		lipgloss.NewStyle().Foreground(theme.TextDim).Render(foot),
	)
}

func (m Model) renderArt(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	intensity := []int{1, 2, 3, 4}[min(m.artFrame, 3)]
	lines := make([]string, height)
	for y := 0; y < height; y++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			ch := ' '
			switch m.view {
			case viewBrowse:
				if (x+y)%9 < intensity {
					ch = '/'
				}
				if (x*2+y)%15 == 0 && intensity > 2 {
					ch = 'o'
				}
			case viewProfile:
				if y >= height-((x/3)%max(2, height-1))-1 && x%3 == 0 && intensity > 1 {
					ch = '|'
				}
			case viewRandom:
				cx, cy := width/2, height/2
				d := abs(x-cx) + abs(y-cy)
				if d%8 < intensity-1 {
					ch = 'o'
				}
			case viewSettings:
				if y%2 == 0 && x > 1 && x < width-2 && intensity > 1 {
					ch = '-'
				}
				if (x+y)%13 == 0 && intensity > 2 {
					ch = 'o'
				}
			default:
				if (x+y)%7 < intensity-1 {
					ch = '#'
				}
			}
			b.WriteRune(ch)
		}
		lines[y] = b.String()
	}
	return lipgloss.NewStyle().Foreground(theme.Surface2).Render(strings.Join(lines, "\n"))
}

func (m *Model) refreshCollections() {
	tot := map[string]int{}
	sol := map[string]int{}
	seen := map[string]bool{}
	for _, e := range m.index.Challenges {
		tot[e.Category]++
	}
	for _, s := range m.profile.Solves {
		if !s.Passed || seen[s.ChallengeID] {
			continue
		}
		e, ok := m.findEntry(s.ChallengeID)
		if !ok {
			continue
		}
		seen[s.ChallengeID] = true
		sol[e.Category]++
	}
	keys := make([]string, 0, len(tot))
	for k := range tot {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	m.collections = m.collections[:0]
	for _, k := range keys {
		m.collections = append(m.collections, collection{Name: k, Solved: sol[k], Total: tot[k]})
	}
}

func (m Model) collectionEntries() []challenge.IndexEntry {
	out := []challenge.IndexEntry{}
	for _, e := range m.index.Challenges {
		if e.Category == m.activeCollect {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}

func (m Model) randomEntries() []challenge.IndexEntry {
	if m.randomFilter == 0 {
		return m.index.Challenges
	}
	var want challenge.Difficulty = challenge.Easy
	if m.randomFilter == 2 {
		want = challenge.Medium
	}
	if m.randomFilter == 3 {
		want = challenge.Hard
	}
	out := []challenge.IndexEntry{}
	for _, e := range m.index.Challenges {
		if e.Difficulty == want {
			out = append(out, e)
		}
	}
	return out
}

func (m *Model) resetSettings() {
	m.settings = []setting{
		newSetting("AI Provider", m.cfg.AIProvider, []string{"claude", "openai", "none"}),
		newSetting("Default Language", m.cfg.PreferredLanguage, []string{"go", "python", "javascript", "cpp", "rust"}),
		newSetting("Theme", m.cfg.Theme, []string{"dark"}),
		{Key: "Repo Path", Value: m.cfg.RepoPath},
		newSetting("Editor Tab Width", strconv.Itoa(m.cfg.TabWidth), []string{"2", "4"}),
	}
	if m.settingCur >= len(m.settings) {
		m.settingCur = max(0, len(m.settings)-1)
	}
}

func newSetting(key, val string, options []string) setting {
	idx := 0
	for i, o := range options {
		if o == val {
			idx = i
			break
		}
	}
	return setting{Key: key, Value: val, Options: options, Index: idx}
}

func (m *Model) applySetting(s *setting) {
	if s == nil || len(s.Options) == 0 {
		return
	}
	s.Value = s.Options[s.Index]
	switch s.Key {
	case "AI Provider":
		m.cfg.AIProvider = s.Value
	case "Default Language":
		m.cfg.PreferredLanguage = s.Value
	case "Theme":
		m.cfg.Theme = s.Value
	case "Editor Tab Width":
		if w, err := strconv.Atoi(s.Value); err == nil {
			m.cfg.TabWidth = w
		}
	}
}

func (m Model) mode() string {
	switch m.view {
	case viewBrowse:
		return "BROWSE"
	case viewProfile:
		return "PROFILE"
	case viewRandom:
		return "RANDOM"
	case viewSettings:
		return "SETTINGS"
	default:
		return "HOME"
	}
}

func (m Model) hints() []components.KeyHint {
	if m.view == viewBrowse && m.browseOpen {
		return []components.KeyHint{{Key: "esc", Description: "back"}, {Key: "j/k", Description: "navigate"}, {Key: "enter", Description: "start"}, {Key: "q", Description: "quit"}}
	}
	if m.view == viewRandom {
		return []components.KeyHint{{Key: "left/right", Description: "filter"}, {Key: "enter", Description: "start"}, {Key: "q", Description: "quit"}}
	}
	if m.view == viewSettings {
		if m.editingSetting {
			return []components.KeyHint{{Key: "left/right", Description: "choose"}, {Key: "enter", Description: "apply"}, {Key: "esc", Description: "cancel"}}
		}
		return []components.KeyHint{{Key: "j/k", Description: "navigate"}, {Key: "enter", Description: "edit"}, {Key: "q", Description: "quit"}}
	}
	return []components.KeyHint{{Key: "j/k", Description: "navigate"}, {Key: "enter", Description: "select"}, {Key: "q", Description: "quit"}}
}

func (m Model) sidebarWidth() int {
	w := m.width / 4
	if w < 24 {
		w = 24
	}
	if w > 30 {
		w = 30
	}
	if m.width-w < 30 {
		w = max(18, m.width-30)
	}
	if w >= m.width {
		w = max(1, m.width-1)
	}
	return w
}

func shortcut(k string) bool {
	return k == "d" || k == "b" || k == "p" || k == "r" || k == "s" || k == "q"
}
func animateArtTick() tea.Cmd {
	return tea.Tick(90*time.Millisecond, func(time.Time) tea.Msg { return artTickMsg{} })
}
func startChallenge(e challenge.IndexEntry) tea.Cmd {
	return func() tea.Msg { return SelectChallengeMsg{Entry: e} }
}
func configChanged(cfg config.Config) tea.Cmd {
	return func() tea.Msg { return ConfigChangedMsg{Config: cfg} }
}

func (m Model) recentRows(width int) string {
	rec := m.profile.RecentSolves(5)
	if len(rec) == 0 {
		return lipgloss.NewStyle().Foreground(theme.TextDim).Render("No sessions yet")
	}
	rows := []string{}
	for i := len(rec) - 1; i >= 0; i-- {
		r := rec[i]
		ent, ok := m.findEntry(r.ChallengeID)
		t := r.ChallengeID
		d := ""
		if ok {
			t, d = ent.Title, string(ent.Difficulty)
		}
		mark := lipgloss.NewStyle().Foreground(theme.Green).Render("v")
		if !r.Passed {
			mark = lipgloss.NewStyle().Foreground(theme.TextDim).Render(".")
		}
		badge := ""
		if d != "" {
			badge = " " + theme.DifficultyStyle(d).Render(d)
		}
		rows = append(rows, mark+" "+cutToWidth(t, max(12, width-28))+badge+" "+lipgloss.NewStyle().Foreground(theme.TextDim).Render(ago(r.Timestamp)))
	}
	return strings.Join(rows, "\n")
}

func (m Model) collectionGrid(width int) string {
	cols := 1
	if width >= 76 {
		cols = 2
	}
	gap := 2
	cw := width
	if cols > 1 {
		cw = (width - gap) / 2
	}
	cards := []string{}
	for i, c := range m.collections {
		cards = append(cards, collectionCard(c, cw, i == m.browseCursor))
	}
	rows := []string{}
	for i := 0; i < len(cards); i += cols {
		rows = append(rows, strings.Join(cards[i:min(len(cards), i+cols)], strings.Repeat(" ", gap)))
	}
	return strings.Join(rows, "\n")
}

func (m Model) challengeRow(e challenge.IndexEntry, width int, cur bool) string {
	mark := lipgloss.NewStyle().Foreground(theme.TextDim).Render(".")
	if m.profile.IsSolved(e.ID) {
		mark = lipgloss.NewStyle().Foreground(theme.Green).Render("v")
	}
	pref := "  "
	st := lipgloss.NewStyle().Foreground(theme.TextDim)
	if cur {
		pref = "> "
		st = st.Foreground(theme.TextPrimary).Background(theme.Surface)
	}
	tags := ""
	if len(e.Tags) > 0 {
		tags = " " + lipgloss.NewStyle().Foreground(theme.TextFaint).Render("#"+strings.Join(e.Tags[:min(2, len(e.Tags))], " #"))
	}
	row := pref + mark + " " + cutToWidth(e.Title, max(12, width-30)) + " " + theme.DifficultyStyle(string(e.Difficulty)).Render(string(e.Difficulty)) + tags
	return st.Width(width).Render(row)
}

func (m Model) settingRow(s setting, width int, cur bool) string {
	label := lipgloss.NewStyle().Foreground(theme.TextDim).Render(s.Key)
	value := lipgloss.NewStyle().Foreground(theme.TextPrimary).Render(s.Value)
	if len(s.Options) == 0 {
		value = lipgloss.NewStyle().Foreground(theme.TextDim).Render(s.Value)
	}
	if cur && m.editingSetting && len(s.Options) > 0 {
		parts := []string{}
		for i, o := range s.Options {
			st := lipgloss.NewStyle().Foreground(theme.TextDim).Padding(0, 1)
			if i == s.Index {
				st = st.Foreground(theme.TextPrimary).Background(theme.Surface2)
			}
			parts = append(parts, st.Render(o))
		}
		value = strings.Join(parts, " ")
	}
	line := label + " "
	gap := width - lipgloss.Width(line) - lipgloss.Width(value)
	if gap < 1 {
		gap = 1
	}
	line += strings.Repeat(" ", gap) + value
	st := lipgloss.NewStyle().Width(width).Padding(0, 1)
	if cur {
		st = st.Background(theme.Surface)
	}
	return st.Render(line)
}

func (m Model) languagePills() string {
	counts := map[string]int{}
	for _, s := range m.profile.Solves {
		if s.Passed {
			counts[s.Language]++
		}
	}
	if len(counts) == 0 {
		counts[m.cfg.PreferredLanguage] = 0
	}
	langs := make([]string, 0, len(counts))
	for l := range counts {
		langs = append(langs, l)
	}
	sort.Strings(langs)
	pills := []string{}
	for _, l := range langs {
		st := lipgloss.NewStyle().Foreground(theme.TextDim).Background(theme.Surface2).Padding(0, 1)
		if l == m.cfg.PreferredLanguage {
			st = st.Foreground(theme.TextPrimary).Border(lipgloss.NormalBorder()).BorderForeground(theme.Red)
		}
		pills = append(pills, st.Render(l+" "+strconv.Itoa(counts[l])))
	}
	return strings.Join(pills, " ")
}

func (m Model) diffCounts() (int, int, int) {
	s := map[string]bool{}
	for _, r := range m.profile.Solves {
		if r.Passed {
			s[r.ChallengeID] = true
		}
	}
	e, md, h := 0, 0, 0
	for _, c := range m.index.Challenges {
		if !s[c.ID] {
			continue
		}
		switch c.Difficulty {
		case challenge.Easy:
			e++
		case challenge.Medium:
			md++
		case challenge.Hard:
			h++
		}
	}
	return e, md, h
}

func (m Model) categoryBars(width int) string {
	total, solved, seen := map[string]int{}, map[string]int{}, map[string]bool{}
	for _, c := range m.index.Challenges {
		total[c.Category]++
	}
	for _, r := range m.profile.Solves {
		if !r.Passed || seen[r.ChallengeID] {
			continue
		}
		c, ok := m.findEntry(r.ChallengeID)
		if !ok {
			continue
		}
		seen[r.ChallengeID] = true
		solved[c.Category]++
	}
	cats := make([]string, 0, len(total))
	for c := range total {
		cats = append(cats, c)
	}
	sort.Strings(cats)
	rows := []string{}
	for _, c := range cats {
		rows = append(rows,
			lipgloss.NewStyle().Foreground(theme.TextDim).Render(cutToWidth(c, max(8, width/3)))+" "+
				progress(max(10, width/3), solved[c], total[c])+" "+
				lipgloss.NewStyle().Foreground(theme.TextFaint).Render(strconv.Itoa(solved[c])+"/"+strconv.Itoa(total[c])),
		)
	}
	if len(rows) == 0 {
		return lipgloss.NewStyle().Foreground(theme.TextDim).Render("No category data")
	}
	return strings.Join(rows, "\n")
}

func (m Model) avgSolve() string {
	total, count := time.Duration(0), 0
	for _, s := range m.profile.Solves {
		if s.Passed {
			total += s.Duration.Duration
			count++
		}
	}
	if count == 0 {
		return "--"
	}
	a := total / time.Duration(count)
	if a < time.Minute {
		return strconv.Itoa(int(a.Seconds())) + "s"
	}
	return strconv.Itoa(int(a.Minutes())) + "m " + strconv.Itoa(int(a.Seconds())%60) + "s"
}

func (m Model) findEntry(id string) (challenge.IndexEntry, bool) {
	for _, e := range m.index.Challenges {
		if e.ID == id {
			return e, true
		}
	}
	return challenge.IndexEntry{}, false
}

func statCard(label, value string, color lipgloss.Color, width int) string {
	return lipgloss.NewStyle().Width(width).Background(theme.Surface).Border(lipgloss.NormalBorder()).BorderForeground(theme.Border).Padding(0, 1).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Foreground(theme.TextDim).Render(label),
			lipgloss.NewStyle().Foreground(color).Bold(true).Render(value),
		),
	)
}

func collectionCard(c collection, width int, cur bool) string {
	st := lipgloss.NewStyle().Width(width).Padding(0, 1).Border(lipgloss.RoundedBorder()).BorderForeground(theme.Border)
	if cur {
		st = st.Background(theme.Surface).BorderForeground(theme.Red)
	}
	return st.Render(lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true).Render(c.Name),
		lipgloss.NewStyle().Foreground(theme.TextDim).Render(strconv.Itoa(c.Solved)+" / "+strconv.Itoa(c.Total)),
		progress(max(8, width-8), c.Solved, c.Total),
	))
}

func progress(width, solved, total int) string {
	if width <= 0 {
		return ""
	}
	if total <= 0 {
		total = 1
	}
	fill := solved * width / total
	if fill > width {
		fill = width
	}
	return lipgloss.NewStyle().Foreground(theme.Red).Render(strings.Repeat("=", fill)) +
		lipgloss.NewStyle().Foreground(theme.Surface2).Render(strings.Repeat("=", width-fill))
}

func formatStreak(days int) string {
	if days == 1 {
		return "1 day"
	}
	return strconv.Itoa(days) + " days"
}
func ago(t time.Time) string {
	if t.IsZero() {
		return "just now"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return strconv.Itoa(int(d.Minutes())) + "m ago"
	case d < 24*time.Hour:
		return strconv.Itoa(int(d.Hours())) + "h ago"
	case d < 7*24*time.Hour:
		return strconv.Itoa(int(d.Hours()/24)) + "d ago"
	default:
		return t.Format("2006-01-02")
	}
}

func trimHeight(s string, h int) string {
	if h <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) <= h {
		return s
	}
	return strings.Join(lines[:h], "\n")
}

func cutToWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= width {
		return s
	}
	runes := []rune(s)
	for len(runes) > 0 && lipgloss.Width(string(runes)) > width {
		runes = runes[:len(runes)-1]
	}
	return string(runes)
}
func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
