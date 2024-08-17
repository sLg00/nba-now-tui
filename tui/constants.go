package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
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
}

// Date is a helper function which returns the current date in a pre-specified format
func Date() string {
	date := time.Now().Format("2006-01-02")
	return date
}
