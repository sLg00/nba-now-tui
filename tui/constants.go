package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	Program    *tea.Program
	WindowSize tea.WindowSizeMsg
)

type errMsg struct{ error }

type keymap struct {
	Back  key.Binding
	Quit  key.Binding
	Enter key.Binding
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Tab   key.Binding
	Space key.Binding
}

var DocStyle = lipgloss.NewStyle().Margin(0, 2)

var TableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color("240"))

// HelpStyle styling for help context menu
var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

var Keymap = keymap{
	Back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "back")),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit")),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select")),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "up")),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "down")),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("left", "left")),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("right", "right")),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch")),
	Space: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "mark for selection")),
}

// HelpFooter returns a fully-fledged help footer to be used in all nba-tui views
func HelpFooter() string {
	localMap := []struct{ Key, Desc string }{
		{Keymap.Back.Help().Key, Keymap.Back.Help().Desc},
		{Keymap.Quit.Help().Key, Keymap.Quit.Help().Desc},
		{Keymap.Enter.Help().Key, Keymap.Enter.Help().Desc},
		{Keymap.Space.Help().Key, Keymap.Space.Help().Desc}}

	helpItems := localMap

	var builder strings.Builder
	for i, item := range helpItems {
		helpItem := item.Key + ": " + item.Desc
		builder.WriteString(helpItem)

		if i < len(helpItems)-1 {
			builder.WriteString(" | ")
		}
	}

	helpMenu := strings.TrimSpace(builder.String())
	return helpMenu
}
