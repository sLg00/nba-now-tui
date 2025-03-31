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

var DocStyle = lipgloss.NewStyle().Margin(2, 2).BorderStyle(lipgloss.HiddenBorder())

var TableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder()).
	Align(lipgloss.Center, lipgloss.Center).
	BorderForeground(lipgloss.Color("240"))

var InvisibleTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.HiddenBorder()).
	Align(lipgloss.Center, lipgloss.Center).
	BorderForeground(lipgloss.Color("#000000"))

var ViewPortBaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("4")).
	Margin(1, 1).Padding(1, 1)

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

// CenterStyle takes a variable width and returns a centered style based on that. Used to align content in viewports
func CenterStyle(w int) lipgloss.Style {
	return lipgloss.NewStyle().Width(w).Align(lipgloss.Center)
}

// TeamColor returns a lipgloss.Color object, representing the dominant color for the specific team
func TeamColor(s string) lipgloss.Color {

	//this map contains the names:colors mappings
	clrs := map[string]lipgloss.Color{
		"76ers": lipgloss.Color("#ED174C"),
		"Bucks": lipgloss.Color("#00471B"),
	}
	if clr, ok := clrs[s]; ok {
		return clr
	}
	return lipgloss.Color("#FFFFFF")
}

func TeamViewPortStyle(clr lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(clr).
		Margin(1, 1).Padding(1, 1)
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
