package leaderboard

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/config"
	"github.com/spar-cli/spar/internal/friends"
	"github.com/spar-cli/spar/internal/profile"
	"github.com/spar-cli/spar/internal/rank"
	"github.com/spar-cli/spar/internal/ui/components"
	"github.com/spar-cli/spar/internal/ui/theme"
)

type SyncRequestMsg struct{}

type subView int

const (
	listView subView = iota
	detailView
	addingFriend
)

type leaderboardEntry struct {
	Username   string
	TotalSP    int
	Rank       string
	Division   int
	Streak     int
	IsYou      bool
	HasProfile bool
	Profile    *friends.PublicProfile
}

type Model struct {
	width         int
	height        int
	profile       *profile.Profile
	cfg           config.Config
	entries       []leaderboardEntry
	noProfile     []leaderboardEntry
	cursor        int
	view          subView
	syncing       bool
	lastSync      time.Time
	friendList    []friends.Friend
	addInput      string
	addResult     string
	confirmRemove bool
}

func New(p *profile.Profile, cfg config.Config) Model {
	m := Model{
		profile: p,
		cfg:     cfg,
	}
	m.loadFriendData()
	return m
}

func (m Model) WithFriendData(results []friends.SyncResult, meta friends.SyncMeta) Model {
	m.syncing = false
	m.lastSync = meta.LastSync
	m.loadFriendData()
	return m
}

func (m Model) InSubView() bool {
	return m.view == detailView || m.view == addingFriend || m.confirmRemove
}

func (m *Model) loadFriendData() {
	fl, _ := friends.LoadFriends(config.FriendsFilePath())
	m.friendList = fl

	var ranked []leaderboardEntry
	var noProf []leaderboardEntry

	selfEntry := leaderboardEntry{
		Username:   m.profile.Username,
		TotalSP:    m.profile.TotalSP,
		Rank:       m.profile.CurrentTier,
		Division:   m.profile.CurrentDivision,
		Streak:     m.profile.Streak,
		IsYou:      true,
		HasProfile: true,
	}
	if selfEntry.Username == "" {
		selfEntry.Username = "you"
	}
	ranked = append(ranked, selfEntry)

	for _, f := range fl {
		cached, err := friends.LoadCached(config.DataDir(), f.Username)
		if err != nil {
			noProf = append(noProf, leaderboardEntry{
				Username:   f.Username,
				HasProfile: false,
			})
			continue
		}
		p := cached
		ranked = append(ranked, leaderboardEntry{
			Username:   f.Username,
			TotalSP:    cached.TotalSP,
			Rank:       cached.Rank,
			Division:   cached.Division,
			Streak:     cached.Streak,
			IsYou:      false,
			HasProfile: true,
			Profile:    &p,
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].TotalSP > ranked[j].TotalSP
	})

	m.entries = ranked
	m.noProfile = noProf
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.view == addingFriend {
			return m.updateAddInput(msg)
		}
		if m.confirmRemove {
			return m.updateConfirmRemove(msg)
		}
		if m.view == detailView {
			return m.updateDetail(msg)
		}
		return m.updateList(msg)
	}
	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (Model, tea.Cmd) {
	total := len(m.entries) + len(m.noProfile)
	switch msg.String() {
	case "j", "down":
		if total > 0 && m.cursor < total-1 {
			m.cursor++
		}
		return m, nil
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil
	case "enter":
		if m.cursor < len(m.entries) {
			e := m.entries[m.cursor]
			if e.HasProfile && !e.IsYou {
				m.view = detailView
			}
		}
		return m, nil
	case "s":
		m.syncing = true
		return m, requestSync
	case "a":
		m.view = addingFriend
		m.addInput = ""
		m.addResult = ""
		return m, nil
	case "x":
		if m.cursor < len(m.entries) {
			e := m.entries[m.cursor]
			if !e.IsYou {
				m.confirmRemove = true
			}
		} else if m.cursor-len(m.entries) < len(m.noProfile) {
			m.confirmRemove = true
		}
		return m, nil
	}
	return m, nil
}

func (m Model) updateDetail(msg tea.KeyMsg) (Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.view = listView
	}
	return m, nil
}

func (m Model) updateAddInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.view = listView
		return m, nil
	case "enter":
		if m.addInput == "" {
			m.view = listView
			return m, nil
		}
		selfRemote := ""
		if m.cfg.RepoPath != "" {
			selfRemote = gitRemoteURLSafe(m.cfg.RepoPath, m.cfg.GitHub.ForkRemote)
		}
		_, err := friends.AddFriend(config.FriendsFilePath(), m.addInput, selfRemote)
		if err != nil {
			m.addResult = err.Error()
			return m, nil
		}
		m.loadFriendData()
		m.view = listView
		return m, requestSync
	case "backspace":
		if len(m.addInput) > 0 {
			m.addInput = m.addInput[:len(m.addInput)-1]
		}
		return m, nil
	default:
		if len(msg.String()) == 1 || msg.String() == " " {
			m.addInput += msg.String()
		}
		return m, nil
	}
}

func (m Model) updateConfirmRemove(msg tea.KeyMsg) (Model, tea.Cmd) {
	m.confirmRemove = false
	if msg.String() != "y" {
		return m, nil
	}

	var username string
	if m.cursor < len(m.entries) {
		username = m.entries[m.cursor].Username
	} else {
		idx := m.cursor - len(m.entries)
		if idx < len(m.noProfile) {
			username = m.noProfile[idx].Username
		}
	}
	if username != "" {
		_ = friends.RemoveFriend(config.FriendsFilePath(), username)
		m.loadFriendData()
		if m.cursor >= len(m.entries)+len(m.noProfile) && m.cursor > 0 {
			m.cursor--
		}
	}
	return m, nil
}

func (m Model) RenderContent(width, height int) string {
	m.width = width
	m.height = height

	switch m.view {
	case detailView:
		return m.renderDetail()
	default:
		return m.renderList()
	}
}

func (m Model) CurrentHints() []components.KeyHint {
	base := []components.KeyHint{{Key: "tab", Description: "switch focus"}}
	if m.view == addingFriend {
		return append(base,
			components.KeyHint{Key: "enter", Description: "add"},
			components.KeyHint{Key: "esc", Description: "cancel"},
		)
	}
	if m.view == detailView {
		return append(base,
			components.KeyHint{Key: "esc", Description: "back"},
		)
	}
	if m.confirmRemove {
		return append(base,
			components.KeyHint{Key: "y", Description: "confirm"},
			components.KeyHint{Key: "n", Description: "cancel"},
		)
	}
	return append(base,
		components.KeyHint{Key: "j/k", Description: "navigate"},
		components.KeyHint{Key: "enter", Description: "view"},
		components.KeyHint{Key: "s", Description: "sync"},
		components.KeyHint{Key: "a", Description: "add friend"},
		components.KeyHint{Key: "x", Description: "remove"},
		components.KeyHint{Key: "esc", Description: "back"},
	)
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

func requestSync() tea.Msg { return SyncRequestMsg{} }

func (m Model) renderList() string {
	titleStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(theme.TextDim)

	syncText := ""
	if m.syncing {
		syncText = dimStyle.Render("syncing...")
	} else if !m.lastSync.IsZero() {
		syncText = dimStyle.Render("last synced: " + timeAgo(m.lastSync))
	}

	headerLeft := titleStyle.Render("LEADERBOARD")
	gap := max(1, m.width-lipgloss.Width(headerLeft)-lipgloss.Width(syncText)-4)
	header := headerLeft + strings.Repeat(" ", gap) + syncText

	var rows []string
	rows = append(rows, header)
	rows = append(rows, "")

	for i, e := range m.entries {
		rows = append(rows, m.renderEntry(i, e))
	}

	if len(m.friendList) == 0 && m.view != addingFriend {
		rows = append(rows, "")
		emptyStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
		rows = append(rows, emptyStyle.Render("No friends yet. Press  a  to add a friend by GitHub username."))
	}

	if len(m.noProfile) > 0 {
		rows = append(rows, "")
		sepStyle := lipgloss.NewStyle().Foreground(theme.TextFaint)
		rows = append(rows, sepStyle.Render("── friends with no profile ──"))
		for i, e := range m.noProfile {
			idx := len(m.entries) + i
			rows = append(rows, m.renderNoProfileEntry(idx, e))
		}
	}

	if m.view == addingFriend {
		rows = append(rows, "")
		inputStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary)
		prompt := inputStyle.Render("Add friend: " + m.addInput + "█")
		rows = append(rows, prompt)
		if m.addResult != "" {
			errStyle := lipgloss.NewStyle().Foreground(theme.Red)
			rows = append(rows, errStyle.Render("  "+m.addResult))
		}
	}

	if m.confirmRemove {
		rows = append(rows, "")
		warnStyle := lipgloss.NewStyle().Foreground(theme.Amber)
		rows = append(rows, warnStyle.Render("Remove this friend? (y/n)"))
	}

	return strings.Join(rows, "\n")
}

func (m Model) renderEntry(idx int, e leaderboardEntry) string {
	isCursor := idx == m.cursor
	rankNum := fmt.Sprintf("#%d", idx+1)

	numStyle := lipgloss.NewStyle().Foreground(theme.TextDim).Width(5).Align(lipgloss.Right)
	nameStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary)
	spStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary).Width(12).Align(lipgloss.Right)

	ri := rank.Calculate(e.TotalSP)
	icon := rank.RenderInline(ri.TierIndex)
	rankColor := lipgloss.Color(ri.Tier.Color)
	rankLabel := lipgloss.NewStyle().Foreground(rankColor).Render(
		ri.Tier.Name + " " + rank.DivisionLabel(ri.Division))

	name := e.Username
	if e.IsYou {
		name += lipgloss.NewStyle().Foreground(theme.TextDim).Render(" (you)")
	}

	streak := lipgloss.NewStyle().Foreground(theme.TextDim).Render("  —")
	if e.Streak > 0 {
		streak = lipgloss.NewStyle().Foreground(theme.Amber).Render(
			fmt.Sprintf("  🔥 %dd", e.Streak))
	}

	line := numStyle.Render(rankNum) + "  " + icon + "  " + nameStyle.Render(name) +
		"  " + rankLabel + spStyle.Render(formatSP(e.TotalSP)+" SP") + streak

	if isCursor {
		return lipgloss.NewStyle().Background(theme.Surface2).
			Width(max(1, m.width-4)).Render(line)
	}
	return line
}

func (m Model) renderNoProfileEntry(idx int, e leaderboardEntry) string {
	isCursor := idx == m.cursor
	dimStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
	nameStyle := lipgloss.NewStyle().Foreground(theme.TextDim)

	line := "       " + nameStyle.Render(e.Username) + "  " + dimStyle.Render("not published")

	if isCursor {
		return lipgloss.NewStyle().Background(theme.Surface2).
			Width(max(1, m.width-4)).Render(line)
	}
	return line
}

func (m Model) renderDetail() string {
	if m.cursor >= len(m.entries) {
		return ""
	}
	e := m.entries[m.cursor]
	if e.Profile == nil {
		return ""
	}
	p := e.Profile
	ri := rank.Calculate(p.TotalSP)
	rankColor := lipgloss.Color(ri.Tier.Color)

	titleStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true)
	rankStyle := lipgloss.NewStyle().Foreground(rankColor)
	dimStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
	labelStyle := lipgloss.NewStyle().Foreground(theme.TextMid).Bold(true)

	var rows []string

	header := titleStyle.Render(p.Username) + " · " +
		rankStyle.Render(ri.Tier.Name+" "+rank.DivisionLabel(ri.Division)) + " · " +
		titleStyle.Render(formatSP(p.TotalSP)+" SP")
	rows = append(rows, header)
	rows = append(rows, "")

	icon := rank.RenderIcon(ri.TierIndex)
	progress := renderProgressBar(ri, rankColor, 24)
	rows = append(rows, icon)
	rows = append(rows, rankStyle.Render(rank.FullName(ri))+"  "+progress+" "+
		dimStyle.Render(fmt.Sprintf("%d / %d SP", p.TotalSP, ri.NextSP)))
	if !ri.IsMax {
		rows = append(rows, dimStyle.Render(fmt.Sprintf("next: %s", ri.NextName)))
	}
	rows = append(rows, "")

	if len(p.TrackMedals) > 0 {
		rows = append(rows, labelStyle.Render("TRACKS"))
		var trackLine []string
		for _, name := range rank.TrackNames {
			key := trackKey(name)
			medal := p.TrackMedals[key]
			ic := medalIcon(medal)
			trackLine = append(trackLine, ic+" "+shortTrackName(name))
		}
		for i := 0; i < len(trackLine); i += 3 {
			end := i + 3
			if end > len(trackLine) {
				end = len(trackLine)
			}
			row := ""
			for j := i; j < end; j++ {
				row += lipgloss.NewStyle().Width(22).Render(trackLine[j])
			}
			rows = append(rows, row)
		}
		rows = append(rows, "")
	}

	if len(p.Languages) > 0 {
		rows = append(rows, labelStyle.Render("LANGUAGES"))
		var langParts []string
		for lang, count := range p.Languages {
			langParts = append(langParts, fmt.Sprintf("%s (%d)", lang, count))
		}
		rows = append(rows, dimStyle.Render(strings.Join(langParts, "  ")))
		rows = append(rows, "")
	}

	rows = append(rows, labelStyle.Render("SOLVES"))
	easy := p.Solves["easy"]
	medium := p.Solves["medium"]
	hard := p.Solves["hard"]
	rows = append(rows, fmt.Sprintf("Easy %d   Medium %d   Hard %d   Total %d / %d",
		easy, medium, hard, p.TotalSolved, p.TotalChallenges))
	rows = append(rows, "")

	if !p.LastUpdated.IsZero() {
		rows = append(rows, dimStyle.Render("last updated: "+timeAgo(p.LastUpdated)))
	}

	return strings.Join(rows, "\n")
}

func renderProgressBar(ri rank.RankInfo, color lipgloss.Color, width int) string {
	filled := int(ri.Progress * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(theme.TextFaint).Render(strings.Repeat("░", width-filled))
}

func RenderArt(width, height int) string {
	if height <= 0 || width <= 0 {
		return ""
	}
	artStyle := lipgloss.NewStyle().Foreground(theme.TextFaint)

	patterns := []string{
		"      ╔═══╗",
		"    ╔═╣ 1 ╠═╗",
		"    ║ ╚═══╝ ║",
		"  ╔═╣  2  3  ╠═╗",
		"  ║ ╚═══════╝ ║",
		"  ╚═══════════╝",
	}

	var lines []string
	for i := 0; i < height && i < len(patterns); i++ {
		pad := (width - lipgloss.Width(patterns[i])) / 2
		if pad < 0 {
			pad = 0
		}
		lines = append(lines, strings.Repeat(" ", pad)+patterns[i])
	}

	return artStyle.Render(strings.Join(lines, "\n"))
}

func formatSP(sp int) string {
	if sp < 1000 {
		return fmt.Sprintf("%d", sp)
	}
	return fmt.Sprintf("%d,%03d", sp/1000, sp%1000)
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func gitRemoteURLSafe(repoPath, remoteName string) string {
	if repoPath == "" || remoteName == "" {
		return ""
	}
	return ""
}

func medalIcon(medal string) string {
	switch medal {
	case "gold":
		return lipgloss.NewStyle().Foreground(theme.MedalGold).Render("★")
	case "silver":
		return lipgloss.NewStyle().Foreground(theme.MedalSilver).Render("◆")
	case "bronze":
		return lipgloss.NewStyle().Foreground(theme.MedalBronze).Render("△")
	default:
		return lipgloss.NewStyle().Foreground(theme.TextFaint).Render("○")
	}
}

func trackKey(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, "the_", "")
	return s
}

func shortTrackName(name string) string {
	name = strings.TrimPrefix(name, "The ")
	parts := strings.Fields(name)
	if len(parts) > 2 {
		return parts[0] + " " + parts[1]
	}
	return name
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
