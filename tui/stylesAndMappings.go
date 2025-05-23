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
	BorderStyle(lipgloss.NormalBorder()).
	Align(lipgloss.Center, lipgloss.Center).
	BorderForeground(lipgloss.Color("240"))

var InvisibleTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.HiddenBorder()).
	Align(lipgloss.Center, lipgloss.Center).
	BorderForeground(lipgloss.Color("#000000"))

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
		"76ers":        lipgloss.Color("#ED174C"),
		"Bucks":        lipgloss.Color("#00471B"),
		"Bulls":        lipgloss.Color("#CE1141"),
		"Cavaliers":    lipgloss.Color("#860038"),
		"Celtics":      lipgloss.Color("#007A33"),
		"Clippers":     lipgloss.Color("#1d428a"),
		"Grizzlies":    lipgloss.Color("#5D76A9"),
		"Hawks":        lipgloss.Color("#C8102E"),
		"Heat":         lipgloss.Color("#db3eb1"),
		"Hornets":      lipgloss.Color("#00788C"),
		"Jazz":         lipgloss.Color("#F9A01B"),
		"Kings":        lipgloss.Color("#5a2d81"),
		"Knicks":       lipgloss.Color("#F58426"),
		"Lakers":       lipgloss.Color("#552583"),
		"Magic":        lipgloss.Color("#C4ced4"),
		"Mavericks":    lipgloss.Color("#00538C"),
		"Nets":         lipgloss.Color("#FFFFFF"),
		"Nuggets":      lipgloss.Color("#1D428A"),
		"Pacers":       lipgloss.Color("#FDBB30"),
		"Pelicans":     lipgloss.Color("#85714D"),
		"Pistons":      lipgloss.Color("#1d42ba"),
		"Raptors":      lipgloss.Color("#ce1141"),
		"Rockets":      lipgloss.Color("#CE1141"),
		"Spurs":        lipgloss.Color("#c4ced4"),
		"Suns":         lipgloss.Color("#e56020"),
		"Thunder":      lipgloss.Color("#007ac1"),
		"Timberwolves": lipgloss.Color("#78BE20"),
		"TrailBlazers": lipgloss.Color("#E03A3E"),
		"Warriors":     lipgloss.Color("#ffc72c"),
		"Wizards":      lipgloss.Color("#002B5C"),
	}
	if clr, ok := clrs[s]; ok {
		return clr
	}
	//default is white
	return lipgloss.Color("#FFFFFF")
}

// TeamViewPortStyle just returns the style with the corresponding team color
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
