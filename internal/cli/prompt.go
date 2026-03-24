package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/NeriCarcasci/spar/internal/ui/theme"
)

var scanner = bufio.NewScanner(os.Stdin)

var (
	successMark = lipgloss.NewStyle().Foreground(theme.Green).Render("✓")
	failMark    = lipgloss.NewStyle().Foreground(theme.Red).Render("✗")
	labelStyle  = lipgloss.NewStyle().Foreground(theme.TextPrimary)
	dimStyle    = lipgloss.NewStyle().Foreground(theme.TextDim)
	hintStyle   = lipgloss.NewStyle().Foreground(theme.TextMid)
	nameStyle   = lipgloss.NewStyle().Foreground(theme.Red).Bold(true)
)

func PromptChoice(question string, options []string) int {
	fmt.Println(hintStyle.Render(question))
	for i, opt := range options {
		num := lipgloss.NewStyle().Foreground(theme.TextDim).Render(fmt.Sprintf("[%d]", i+1))
		fmt.Printf("  %s %s\n", num, labelStyle.Render(opt))
	}
	for {
		fmt.Print(dimStyle.Render("\n> "))
		if !scanner.Scan() {
			return -1
		}
		input := strings.TrimSpace(scanner.Text())
		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > len(options) {
			fmt.Printf(hintStyle.Render("Enter a number between 1 and %d.")+"\n", len(options))
			continue
		}
		return n
	}
}

func PromptString(question, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s %s\n", hintStyle.Render(question), dimStyle.Render("(default: "+defaultVal+")"))
	} else {
		fmt.Println(hintStyle.Render(question))
	}
	fmt.Print(dimStyle.Render("> "))
	if !scanner.Scan() {
		return defaultVal
	}
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return defaultVal
	}
	return input
}

func PromptSecret(question string) string {
	fmt.Println(hintStyle.Render(question))
	fmt.Print(dimStyle.Render("> "))
	if !scanner.Scan() {
		return ""
	}
	return strings.TrimSpace(scanner.Text())
}

func PromptYesNo(question string, defaultYes bool) bool {
	hint := "Y/n"
	if !defaultYes {
		hint = "y/N"
	}
	fmt.Printf("%s %s ", hintStyle.Render(question), dimStyle.Render("["+hint+"]"))
	if !scanner.Scan() {
		return defaultYes
	}
	input := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if input == "" {
		return defaultYes
	}
	return input == "y" || input == "yes"
}

func PrintSuccess(format string, args ...interface{}) {
	msg := lipgloss.NewStyle().Foreground(theme.Green).Render(fmt.Sprintf(format, args...))
	fmt.Printf("%s %s\n", successMark, msg)
}

func PrintError(format string, args ...interface{}) {
	msg := lipgloss.NewStyle().Foreground(theme.Red).Render(fmt.Sprintf(format, args...))
	fmt.Fprintf(os.Stderr, "%s %s\n", failMark, msg)
}
