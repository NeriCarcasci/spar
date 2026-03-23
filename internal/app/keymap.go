package app

import tea "github.com/charmbracelet/bubbletea"

type KeyMap struct {
	Quit         string
	ForceQuit    string
	Help         string
	Back         string
	BrowseKey    string
	ProfileKey   string
	DashboardKey string
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit:         "q",
		ForceQuit:    "ctrl+c",
		Help:         "?",
		Back:         "esc",
		BrowseKey:    "b",
		ProfileKey:   "p",
		DashboardKey: "d",
	}
}

func isGlobalQuit(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+c"
}
