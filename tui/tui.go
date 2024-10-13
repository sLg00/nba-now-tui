package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

// RenderUI is the entrypoint into the TUI
func RenderUI() {
	m, _ := InitMenu()
	Program = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := Program.Run(); err != nil {
		log.Fatal(err)

	}
}
