package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spar-cli/spar/internal/app"
	"github.com/spar-cli/spar/internal/challenge"
	"github.com/spar-cli/spar/internal/config"
	"github.com/spar-cli/spar/internal/generate"
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
