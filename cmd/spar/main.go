package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/app"
	"github.com/spar-cli/spar/internal/challenge"
	"github.com/spar-cli/spar/internal/config"
	"github.com/spar-cli/spar/internal/friends"
	"github.com/spar-cli/spar/internal/generate"
	"github.com/spar-cli/spar/internal/profile"
	"github.com/spar-cli/spar/internal/rank"
	"github.com/spar-cli/spar/internal/ui/theme"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version":
			printVersion()
			return
		case "generate-index":
			runGenerateIndex()
			return
		case "validate":
			runValidate()
			return
		case "publish":
			runPublish()
			return
		case "friend":
			runFriend()
			return
		case "help", "--help", "-h":
			printUsage()
			return
		default:
			printError("unknown command: %s", os.Args[1])
			fmt.Fprintln(os.Stderr)
			printUsage()
			os.Exit(1)
		}
	}

	runTUI()
}

func runTUI() {
	if err := config.EnsureDirectories(); err != nil {
		printError("creating directories: %v", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		printError("loading config: %v", err)
		os.Exit(1)
	}

	model := app.New(cfg)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		printError("running program: %v", err)
		os.Exit(1)
	}
}

func runGenerateIndex() {
	challengesDir := resolveChallengesDir()
	if err := generate.GenerateIndex(challengesDir); err != nil {
		printError("generating index: %v", err)
		os.Exit(1)
	}
	printSuccess("index.yaml generated")
}

func runValidate() {
	var errors []challenge.ValidationError

	if len(os.Args) > 2 {
		errors = challenge.ValidateChallenge(os.Args[2])
	} else {
		challengesDir := resolveChallengesDir()
		errors = challenge.ValidateAll(challengesDir)
	}

	if len(errors) == 0 {
		printSuccess("all challenges valid")
		return
	}

	for _, e := range errors {
		printError("%s", e.Error())
	}
	os.Exit(1)
}

func resolveChallengesDir() string {
	cfg, err := config.Load()
	if err != nil || cfg.RepoPath == "" {
		cwd, _ := os.Getwd()
		return filepath.Join(cwd, "challenges")
	}
	return filepath.Join(cfg.RepoPath, "challenges")
}

func printVersion() {
	name := cliName().Render("spar")
	fmt.Printf("%s %s\n", name, version)
}

func printUsage() {
	name := cliName()
	sub := cliSubcommand()
	desc := cliDescription()
	arg := cliArg()

	fmt.Println(name.Render("spar"), desc.Render("— code under pressure"))
	fmt.Println()
	fmt.Println(desc.Render("Usage:"))
	fmt.Println("  " + name.Render("spar") + " " + arg.Render("[command]"))
	fmt.Println()
	fmt.Println(desc.Render("Commands:"))
	fmt.Println("  " + sub.Render("start") + "            " + desc.Render("Start a coding session"))
	fmt.Println("  " + sub.Render("browse") + "           " + desc.Render("Browse challenges"))
	fmt.Println("  " + sub.Render("stats") + "            " + desc.Render("View your profile and stats"))
	fmt.Println("  " + sub.Render("publish") + "          " + desc.Render("Publish profile to GitHub"))
	fmt.Println("  " + sub.Render("friend") + "           " + desc.Render("Manage friends (add/remove/list/sync)"))
	fmt.Println("  " + sub.Render("validate") + "         " + desc.Render("Validate challenge structure"))
	fmt.Println("  " + sub.Render("generate-index") + "   " + desc.Render("Regenerate the challenge index"))
	fmt.Println()
	fmt.Println(desc.Render("Flags:"))
	fmt.Println("  " + arg.Render("--help") + "       " + desc.Render("Show help"))
	fmt.Println("  " + arg.Render("--version") + "    " + desc.Render("Show version"))
}

func printError(format string, args ...interface{}) {
	prefix := lipgloss.NewStyle().Foreground(theme.Red).Render("✗")
	msg := lipgloss.NewStyle().Foreground(theme.Red).Render(fmt.Sprintf(format, args...))
	fmt.Fprintf(os.Stderr, "%s %s\n", prefix, msg)
}

func printSuccess(format string, args ...interface{}) {
	prefix := lipgloss.NewStyle().Foreground(theme.Green).Render("✓")
	msg := lipgloss.NewStyle().Foreground(theme.Green).Render(fmt.Sprintf(format, args...))
	fmt.Printf("%s %s\n", prefix, msg)
}

func cliName() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Red).Bold(true)
}

func cliSubcommand() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.RedLight)
}

func cliDescription() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.TextMid)
}

func cliArg() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.TextPrimary)
}

func runPublish() {
	cfg, err := config.Load()
	if err != nil {
		printError("loading config: %v", err)
		os.Exit(1)
	}

	repoPath := cfg.RepoPath
	if repoPath == "" {
		repoPath, _ = os.Getwd()
	}

	if cfg.GitHub.ForkRemote == "" {
		printError("no git remote configured. Set github.fork_remote in config.")
		os.Exit(1)
	}

	prof, err := profile.Load(config.ProfilePath())
	if err != nil {
		printError("loading profile: %v", err)
		os.Exit(1)
	}

	var idx *challenge.Index
	idxPath := filepath.Join(repoPath, "challenges")
	if _, statErr := os.Stat(filepath.Join(idxPath, "index.yaml")); statErr == nil {
		idx, _ = challenge.LoadIndex(repoPath)
	}

	pub := friends.BuildPublicProfile(prof, idx, version)
	if err := friends.Publish(repoPath, cfg.GitHub.ForkRemote, pub); err != nil {
		printError("%v", err)
		os.Exit(1)
	}

	ri := rank.Calculate(prof.TotalSP)
	printSuccess("Published profile (%s %s · %s SP)",
		ri.Tier.Name, rank.DivisionLabel(ri.Division), formatSP(prof.TotalSP))
}

func runFriend() {
	if len(os.Args) < 3 {
		printError("usage: spar friend <add|remove|list|sync> [args]")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "add":
		runFriendAdd()
	case "remove":
		runFriendRemove()
	case "list":
		runFriendList()
	case "sync":
		runFriendSync()
	default:
		printError("unknown friend command: %s", os.Args[2])
		os.Exit(1)
	}
}

func runFriendAdd() {
	if len(os.Args) < 4 {
		printError("usage: spar friend add <url-or-username>")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		printError("loading config: %v", err)
		os.Exit(1)
	}

	selfRemote := ""
	repoPath := cfg.RepoPath
	if repoPath == "" {
		repoPath, _ = os.Getwd()
	}
	if repoPath != "" {
		out, _ := gitOutput(repoPath, "remote", "get-url", cfg.GitHub.ForkRemote)
		selfRemote = strings.TrimSpace(out)
	}

	friend, err := friends.AddFriend(config.FriendsFilePath(), os.Args[3], selfRemote)
	if err != nil {
		printError("%v", err)
		os.Exit(1)
	}

	prof, fetchErr := friends.FetchProfile(friend)
	if fetchErr == nil {
		_ = friends.SaveCached(config.DataDir(), friend.Username, prof)
		ri := rank.Calculate(prof.TotalSP)
		printSuccess("Added %s (%s %s · %s SP)",
			friend.Username, ri.Tier.Name, rank.DivisionLabel(ri.Division), formatSP(prof.TotalSP))
	} else {
		printSuccess("Added %s — no published profile yet", friend.Username)
	}
}

func runFriendRemove() {
	if len(os.Args) < 4 {
		printError("usage: spar friend remove <username>")
		os.Exit(1)
	}

	if err := friends.RemoveFriend(config.FriendsFilePath(), os.Args[3]); err != nil {
		printError("%v", err)
		os.Exit(1)
	}
	printSuccess("Removed %s from friends.", os.Args[3])
}

func runFriendList() {
	fl, err := friends.LoadFriends(config.FriendsFilePath())
	if err != nil {
		printError("loading friends: %v", err)
		os.Exit(1)
	}
	if len(fl) == 0 {
		fmt.Println(lipgloss.NewStyle().Foreground(theme.TextDim).Render("No friends added yet. Use: spar friend add <username>"))
		return
	}

	meta, _ := friends.LoadMeta(config.DataDir())

	nameStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary).Width(14)
	rankStyle := lipgloss.NewStyle().Foreground(theme.TextMid).Width(16)
	spStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary).Width(12).Align(lipgloss.Right)
	dimStyle := lipgloss.NewStyle().Foreground(theme.TextDim)

	for _, f := range fl {
		cached, cErr := friends.LoadCached(config.DataDir(), f.Username)
		if cErr != nil {
			fmt.Printf("%s %s %s\n",
				nameStyle.Render(f.Username),
				rankStyle.Render("—"),
				dimStyle.Render("no profile"))
			continue
		}

		ri := rank.Calculate(cached.TotalSP)
		syncedAgo := ""
		if fm, ok := meta.Results[f.Username]; ok {
			syncedAgo = "synced " + timeAgo(fm.FetchedAt)
		}

		fmt.Printf("%s %s %s   %s\n",
			nameStyle.Render(f.Username),
			rankStyle.Render(ri.Tier.Name+" "+rank.DivisionLabel(ri.Division)),
			spStyle.Render(formatSP(cached.TotalSP)+" SP"),
			dimStyle.Render(syncedAgo))
	}
}

func runFriendSync() {
	fl, err := friends.LoadFriends(config.FriendsFilePath())
	if err != nil {
		printError("loading friends: %v", err)
		os.Exit(1)
	}
	if len(fl) == 0 {
		fmt.Println(lipgloss.NewStyle().Foreground(theme.TextDim).Render("No friends to sync."))
		return
	}

	fmt.Printf("Syncing %d friends...\n", len(fl))
	start := time.Now()
	results := friends.SyncAll(fl)

	meta := friends.SyncMeta{
		LastSync: time.Now().UTC(),
		Results:  make(map[string]friends.FriendMeta),
	}

	okStyle := lipgloss.NewStyle().Foreground(theme.Green)
	errStyle := lipgloss.NewStyle().Foreground(theme.Red)

	for _, r := range results {
		meta.Results[r.Friend.Username] = friends.FriendMeta{
			Status:    r.Status,
			FetchedAt: time.Now().UTC(),
		}

		if r.Profile != nil {
			_ = friends.SaveCached(config.DataDir(), r.Friend.Username, *r.Profile)
			ri := rank.Calculate(r.Profile.TotalSP)
			fmt.Printf("  %s  %s  %s · %s SP\n",
				r.Friend.Username,
				okStyle.Render("✓"),
				ri.Tier.Name+" "+rank.DivisionLabel(ri.Division),
				formatSP(r.Profile.TotalSP))
		} else {
			fmt.Printf("  %s  %s  %s\n",
				r.Friend.Username,
				errStyle.Render("✗"),
				r.Status)
		}
	}

	_ = friends.SaveMeta(config.DataDir(), meta)
	elapsed := time.Since(start)
	fmt.Printf("Synced in %.1fs\n", elapsed.Seconds())
}

func gitOutput(repoPath string, args ...string) (string, error) {
	out, err := execGit(repoPath, args...)
	return string(out), err
}

func execGit(repoPath string, args ...string) ([]byte, error) {
	allArgs := append([]string{"-C", repoPath}, args...)
	return exec.Command("git", allArgs...).Output()
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
		m := int(d.Minutes())
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%dh ago", h)
	default:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
}
