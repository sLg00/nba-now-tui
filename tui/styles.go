package tui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	tableStyle lipgloss.Style
	listStyle  lipgloss.Style
}

func TableStyles() Styles {
	r := lipgloss.DefaultRenderer()
	return Styles{
		tableStyle: r.NewStyle().BorderStyle(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("240")),
	}
}

func ListStyles() Styles {

	return Styles{}
}
