package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/NeriCarcasci/spar/internal/ui/theme"
)

const SidebarQuitID = "quit"

type SidebarItem struct {
	ID        string
	Label     string
	Key       string
	Secondary bool
}

type SidebarEvent struct {
	SelectedChanged bool
	Selected        SidebarItem
	Quit            bool
}

type Sidebar struct {
	items    []SidebarItem
	cursor   int
	selected int
}

func NewSidebar(items []SidebarItem, selectedID string) Sidebar {
	s := Sidebar{items: append([]SidebarItem{}, items...)}
	selected := 0
	for i, item := range s.items {
		if item.ID == selectedID {
			selected = i
			break
		}
	}
	s.cursor = selected
	s.selected = selected
	return s
}

func (s Sidebar) CursorItem() SidebarItem {
	if len(s.items) == 0 {
		return SidebarItem{}
	}
	return s.items[s.cursor]
}

func (s Sidebar) SelectedItem() SidebarItem {
	if len(s.items) == 0 {
		return SidebarItem{}
	}
	return s.items[s.selected]
}

func (s Sidebar) Update(msg tea.Msg) (Sidebar, SidebarEvent) {
	event := SidebarEvent{}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok || len(s.items) == 0 {
		return s, event
	}

	switch keyMsg.String() {
	case "up", "k":
		s.cursor = clampIndex(s.cursor-1, len(s.items))
		return s, event
	case "down", "j":
		s.cursor = clampIndex(s.cursor+1, len(s.items))
		return s, event
	case "enter":
		event = s.selectIndex(s.cursor)
		return s, event
	}

	for i, item := range s.items {
		if item.Key == keyMsg.String() {
			s.cursor = i
			event = s.selectIndex(i)
			return s, event
		}
	}

	return s, event
}

func (s Sidebar) View(width, height int, logo string, focused bool) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	borderColor := theme.Border
	_ = focused
	borderStyle := lipgloss.NewStyle().Foreground(borderColor)
	borderChar := borderStyle.Render("│")

	itemWidth := width - 4
	if itemWidth < 6 {
		itemWidth = 6
	}

	mainItems, secondaryItems := splitItems(s.items)

	logoStyle := lipgloss.NewStyle().
		Foreground(theme.Red).
		Bold(true)

	main := renderItems(mainItems, s, itemWidth)
	secondary := renderItems(secondaryItems, s, itemWidth)

	sep := borderStyle.Render(strings.Repeat("─", itemWidth))

	logoBlock := logoStyle.Render(logo)
	logoLines := lipgloss.Height(logoBlock) + 1
	mainLines := max(1, len(mainItems))
	secondaryLines := 1 + max(1, len(secondaryItems))

	usable := max(1, height-2)
	gap := usable - logoLines - mainLines - secondaryLines
	if gap < 1 {
		gap = 1
	}

	content := strings.Join([]string{
		logoBlock,
		"",
		main,
		blankLines(gap),
		sep,
		secondary,
	}, "\n")

	innerW := width - 4
	if innerW < 4 {
		innerW = 4
	}
	lines := strings.Split(content, "\n")
	out := make([]string, 0, height)

	emptyLine := " " + strings.Repeat(" ", innerW) + "  " + borderChar
	out = append(out, emptyLine)

	for _, line := range lines {
		rendered := " " + lipgloss.NewStyle().Width(innerW).MaxWidth(innerW).Inline(true).Render(line) + "  " + borderChar
		out = append(out, rendered)
	}

	empty := emptyLine
	for len(out) < height {
		out = append(out, empty)
	}
	if len(out) > height {
		out = out[:height]
	}

	return strings.Join(out, "\n")
}

func (s *Sidebar) selectByID(id string) SidebarEvent {
	for i, item := range s.items {
		if item.ID == id {
			s.cursor = i
			return s.selectIndex(i)
		}
	}
	return SidebarEvent{}
}

func (s *Sidebar) selectIndex(index int) SidebarEvent {
	selected := s.items[index]
	event := SidebarEvent{Selected: selected}
	event.SelectedChanged = s.selected != index
	s.selected = index
	if selected.ID == SidebarQuitID {
		event.Quit = true
	}
	return event
}

func splitItems(items []SidebarItem) ([]SidebarItem, []SidebarItem) {
	var mainItems []SidebarItem
	var secondaryItems []SidebarItem
	for _, item := range items {
		if item.Secondary {
			secondaryItems = append(secondaryItems, item)
			continue
		}
		mainItems = append(mainItems, item)
	}
	return mainItems, secondaryItems
}

func renderItems(items []SidebarItem, s Sidebar, width int) string {
	if len(items) == 0 {
		return ""
	}

	lines := make([]string, 0, len(items))
	for _, item := range items {
		index := findItemIndex(s.items, item.ID)
		isCursor := index == s.cursor
		isSelected := index == s.selected
		lines = append(lines, renderItemLine(item, width, isSelected, isCursor))
	}
	return strings.Join(lines, "\n")
}

func renderItemLine(item SidebarItem, width int, isSelected bool, isCursor bool) string {
	indicator := " "
	textStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
	keyStyle := lipgloss.NewStyle().Foreground(theme.TextFaint)
	lineStyle := lipgloss.NewStyle()

	if isSelected {
		indicator = lipgloss.NewStyle().Foreground(theme.Red).Render("›")
		textStyle = textStyle.Foreground(theme.TextPrimary)
		keyStyle = keyStyle.Foreground(theme.TextDim)
		lineStyle = lineStyle.Background(theme.Surface2)
	} else if isCursor {
		textStyle = textStyle.Foreground(theme.TextMid)
		lineStyle = lineStyle.Background(theme.Surface)
	}

	maxLabel := width - 6
	if maxLabel < 4 {
		maxLabel = 4
	}
	label := textStyle.Render(cutToWidth(item.Label, maxLabel))
	key := keyStyle.Render(item.Key)
	left := indicator + " " + label
	spaceCount := width - lipgloss.Width(left) - lipgloss.Width(key)
	if spaceCount < 1 {
		spaceCount = 1
	}

	line := left + strings.Repeat(" ", spaceCount) + key
	return lineStyle.Width(width).Render(line)
}

func findItemIndex(items []SidebarItem, id string) int {
	for i, item := range items {
		if item.ID == id {
			return i
		}
	}
	return 0
}

func clampIndex(i, length int) int {
	if length <= 0 {
		return 0
	}
	if i < 0 {
		return length - 1
	}
	if i >= length {
		return 0
	}
	return i
}

func blankLines(lines int) string {
	if lines <= 0 {
		return ""
	}
	return strings.Repeat("\n", lines-1)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
