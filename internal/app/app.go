package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spar-cli/spar/internal/challenge"
	"github.com/spar-cli/spar/internal/config"
	"github.com/spar-cli/spar/internal/friends"
	"github.com/spar-cli/spar/internal/profile"
	"github.com/spar-cli/spar/internal/rank"
	"github.com/spar-cli/spar/internal/repo"
	"github.com/spar-cli/spar/internal/ui/browser"
	"github.com/spar-cli/spar/internal/ui/dashboard"
	profileview "github.com/spar-cli/spar/internal/ui/profile"
	"github.com/spar-cli/spar/internal/ui/session"
	"github.com/spar-cli/spar/internal/ui/splash"
)

type View int

const (
	SplashView View = iota
	DashboardView
	BrowserView
	SessionView
	ProfileView
)

type Model struct {
	currentView View
	width       int
	height      int
	config      config.Config
	index       *challenge.Index
	profile     *profile.Profile
	syncResult  <-chan repo.SyncResult
	keyMap      KeyMap

	splash    splash.Model
	dashboard dashboard.Model
	browser   browser.Model
	session   session.Model
	profileV  profileview.Model
}

type loadingDoneMsg struct {
	index    *challenge.Index
	profile  *profile.Profile
	repoPath string
}

type syncCompleteMsg struct {
	result repo.SyncResult
}

func New(cfg config.Config) Model {
	return Model{
		currentView: SplashView,
		config:      cfg,
		keyMap:      DefaultKeyMap(),
		splash:      splash.New(),
		index:       &challenge.Index{},
		profile:     &profile.Profile{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.splash.Init(),
		tea.WindowSize(),
		loadData(m.config),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.propagateSize()

	case tea.KeyMsg:
		if isGlobalQuit(msg) {
			return m, tea.Quit
		}

	case loadingDoneMsg:
		m.index = msg.index
		m.profile = msg.profile
		if msg.repoPath != "" {
			m.config.RepoPath = msg.repoPath
		}
		m.dashboard = dashboard.New(m.profile, m.index, m.config)
		m.browser = browser.New(m.index, m.profile)
		m.profileV = profileview.New(m.profile)
		m = m.propagateSize()

		splashModel, cmd := m.splash.Update(splash.LoadingCompleteMsg{})
		m.splash = splashModel
		return m, tea.Batch(cmd, syncFriends(m.config))

	case friendSyncDoneMsg:
		m.dashboard = m.dashboard.HandleFriendSync(msg.results, msg.meta)
		return m, nil

	case syncCompleteMsg:
		if msg.result.Updated && m.config.RepoPath != "" {
			return m, reloadIndex(m.config.RepoPath)
		}

	case reloadedIndexMsg:
		if msg.index != nil {
			m.index = msg.index
			m.browser = browser.New(m.index, m.profile)
			m.dashboard = dashboard.New(m.profile, m.index, m.config)
			m = m.propagateSize()
		}
	}

	var cmd tea.Cmd
	switch m.currentView {
	case SplashView:
		m, cmd = m.updateSplash(msg)
	case DashboardView:
		m, cmd = m.updateDashboard(msg)
	case BrowserView:
		m, cmd = m.updateBrowser(msg)
	case SessionView:
		m, cmd = m.updateSession(msg)
	case ProfileView:
		m, cmd = m.updateProfile(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.currentView {
	case SplashView:
		return m.splash.View()
	case DashboardView:
		return m.dashboard.View()
	case BrowserView:
		return m.browser.View()
	case SessionView:
		return m.session.View()
	case ProfileView:
		return m.profileV.View()
	default:
		return ""
	}
}

func (m Model) updateSplash(msg tea.Msg) (Model, tea.Cmd) {
	splashModel, cmd := m.splash.Update(msg)
	m.splash = splashModel

	if m.splash.ReadyToTransition() {
		m.currentView = DashboardView
	}

	return m, cmd
}

func (m Model) updateDashboard(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case dashboard.SelectChallengeMsg:
		ch, err := challenge.LoadChallenge(m.config.RepoPath, msg.Entry)
		if err != nil {
			return m, nil
		}
		lang := m.config.PreferredLanguage
		code, _ := challenge.LoadSetupCode(ch.Path, lang)
		m.session = session.New(ch, lang).SetSize(m.width, m.height)
		if code != "" {
			m.session = m.session.WithCode(code)
		}
		m.currentView = SessionView
		return m, m.session.Init()
	case dashboard.ConfigChangedMsg:
		m.config = msg.Config
		return m, nil
	case dashboard.ProfileChangedMsg:
		m.profile = msg.Profile
		return m, saveProfile(msg.Profile)
	case dashboard.FriendSyncRequestMsg:
		return m, syncFriends(m.config)
	}

	dashModel, cmd := m.dashboard.Update(msg)
	m.dashboard = dashModel
	return m, cmd
}

func (m Model) updateBrowser(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case browser.NavigateDashboardMsg:
		m.currentView = DashboardView
		return m, nil
	case browser.SelectChallengeMsg:
		ch, err := challenge.LoadChallenge(m.config.RepoPath, msg.Entry)
		if err != nil {
			return m, nil
		}
		lang := m.config.PreferredLanguage
		code, _ := challenge.LoadSetupCode(ch.Path, lang)
		m.session = session.New(ch, lang).SetSize(m.width, m.height)
		if code != "" {
			m.session = m.session.WithCode(code)
		}
		m.currentView = SessionView
		return m, m.session.Init()
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	browserModel, cmd := m.browser.Update(msg)
	m.browser = browserModel
	return m, cmd
}

func (m Model) updateSession(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case session.NavigateBrowserMsg:
		m.currentView = DashboardView
		return m, nil
	case session.SolvedMsg:
		m.handleSolve(msg)
		sessionModel, cmd := m.session.Update(msg)
		m.session = sessionModel
		cmds := []tea.Cmd{cmd, saveProfile(m.profile)}
		if m.config.GitHub.AutoPublish {
			cmds = append(cmds, autoPublish(m.config, m.profile, m.index))
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	sessionModel, cmd := m.session.Update(msg)
	m.session = sessionModel
	return m, cmd
}

func (m *Model) handleSolve(msg session.SolvedMsg) {
	isFirst := m.profile.IsFirstSolve(msg.ChallengeID)
	m.profile.UpdateStreak()

	timeLimitSecs := 0
	for _, e := range m.index.Challenges {
		if e.ID == msg.ChallengeID {
			break
		}
	}

	input := rank.SessionInput{
		Difficulty:    msg.Difficulty,
		DurationSecs:  msg.DurationSecs,
		TimeLimitSecs: timeLimitSecs,
		TestRuns:      msg.TestRuns,
		IsFirstSolve:  isFirst,
		UsedHints:     msg.UsedHints,
	}

	bd := rank.CalculateSP(input, m.profile.Streak)

	record := profile.SolveRecord{
		ChallengeID:  msg.ChallengeID,
		Language:     msg.Language,
		Timestamp:    time.Now(),
		Duration:     profile.Duration{Duration: time.Duration(msg.DurationSecs) * time.Second},
		TestRuns:     msg.TestRuns,
		Passed:       true,
		UsedHints:    msg.UsedHints,
		SPEarned:     bd.TotalSP,
		Multiplier:   bd.Combined,
		IsFirstSolve: isFirst,
	}

	m.profile.Solves = append(m.profile.Solves, record)
	m.profile.TotalSP += bd.TotalSP
	ri := rank.Calculate(m.profile.TotalSP)
	m.profile.CurrentTier = ri.Tier.Name
	m.profile.CurrentDivision = int(ri.Division)
	m.profile.EnsureDefaults()
}

func (m Model) updateProfile(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case profileview.NavigateDashboardMsg:
		m.currentView = DashboardView
		return m, nil
	case tea.KeyMsg:
		if msg.(tea.KeyMsg).String() == "q" {
			return m, tea.Quit
		}
	}

	profileModel, cmd := m.profileV.Update(msg)
	m.profileV = profileModel
	return m, cmd
}



func (m Model) propagateSize() Model {
	m.splash = m.splash.SetSize(m.width, m.height)
	m.dashboard = m.dashboard.SetSize(m.width, m.height)
	m.browser = m.browser.SetSize(m.width, m.height)
	m.session = m.session.SetSize(m.width, m.height)
	m.profileV = m.profileV.SetSize(m.width, m.height)
	return m
}

func loadData(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		config.EnsureDirectories()

		var idx *challenge.Index
		var prof *profile.Profile
		repoPath := resolveRepoPath(cfg.RepoPath)

		if repoPath != "" {
			loaded, err := challenge.LoadIndex(repoPath)
			if err == nil {
				idx = loaded
			}
		}
		if idx == nil {
			idx = &challenge.Index{}
		}

		profPath := config.ProfilePath()
		loaded, err := profile.Load(profPath)
		if err == nil {
			prof = loaded
		}
		if prof == nil {
			prof = &profile.Profile{}
		}

		if _, statErr := os.Stat(profPath); os.IsNotExist(statErr) {
			profile.Save(profPath, prof)
		}

		return loadingDoneMsg{index: idx, profile: prof, repoPath: repoPath}
	}
}

type reloadedIndexMsg struct {
	index *challenge.Index
}

func reloadIndex(repoPath string) tea.Cmd {
	return func() tea.Msg {
		idx, err := challenge.LoadIndex(repoPath)
		if err != nil {
			return reloadedIndexMsg{}
		}
		return reloadedIndexMsg{index: idx}
	}
}

type profileSavedMsg struct{}

func saveProfile(p *profile.Profile) tea.Cmd {
	return func() tea.Msg {
		profile.Save(config.ProfilePath(), p)
		return profileSavedMsg{}
	}
}

var appVersion = "dev"

func SetVersion(v string) { appVersion = v }

type friendSyncDoneMsg struct {
	results []friends.SyncResult
	meta    friends.SyncMeta
}

func autoPublish(cfg config.Config, p *profile.Profile, idx *challenge.Index) tea.Cmd {
	return func() tea.Msg {
		repoPath := cfg.RepoPath
		if repoPath == "" {
			return nil
		}
		pub := friends.BuildPublicProfile(p, idx, appVersion)
		if err := friends.Publish(repoPath, cfg.GitHub.ForkRemote, pub); err != nil {
			log.Printf("auto-publish failed: %v", err)
		}
		return nil
	}
}

func syncFriends(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		fl, err := friends.LoadFriends(config.FriendsFilePath())
		if err != nil || len(fl) == 0 {
			return friendSyncDoneMsg{}
		}
		results := friends.SyncAll(fl)
		meta := friends.SyncMeta{
			LastSync: time.Now().UTC(),
			Results:  make(map[string]friends.FriendMeta),
		}
		for _, r := range results {
			meta.Results[r.Friend.Username] = friends.FriendMeta{
				Status:    r.Status,
				FetchedAt: time.Now().UTC(),
			}
			if r.Profile != nil {
				_ = friends.SaveCached(config.DataDir(), r.Friend.Username, *r.Profile)
			}
		}
		_ = friends.SaveMeta(config.DataDir(), meta)
		return friendSyncDoneMsg{results: results, meta: meta}
	}
}

func resolveRepoPath(configPath string) string {
	if strings.TrimSpace(configPath) != "" {
		return configPath
	}

	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	indexPath := filepath.Join(wd, "challenges", "index.yaml")
	if _, err := os.Stat(indexPath); err == nil {
		return wd
	}

	return ""
}
