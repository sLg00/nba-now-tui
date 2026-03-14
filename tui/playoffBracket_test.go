package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPlayoffBracket_Init_IssuesCmd(t *testing.T) {
	pb, cmd, err := NewPlayoffBracket("2024-25", 0, tea.WindowSizeMsg{Width: 120, Height: 40})
	if err != nil {
		t.Fatalf("NewPlayoffBracket() error: %v", err)
	}
	if pb == nil {
		t.Fatal("NewPlayoffBracket() returned nil model")
	}
	if cmd == nil {
		t.Error("NewPlayoffBracket() returned nil cmd")
	}
}

func TestPlayoffBracket_CursorNavigation(t *testing.T) {
	pb, _, _ := NewPlayoffBracket("2024-25", 0, tea.WindowSizeMsg{Width: 120, Height: 40})

	// Simulate data loaded so keyboard input is unblocked
	loaded, _ := pb.Update(bracketFetchedMsg{})
	ready := loaded.(PlayoffBracket)

	// Move right from East R1 (col 0) to East Semis (col 1)
	model, _ := ready.Update(tea.KeyMsg{Type: tea.KeyRight})
	updated := model.(PlayoffBracket)
	if updated.cursorCol != 1 {
		t.Errorf("after Right from col 0, cursorCol = %d, want 1", updated.cursorCol)
	}

	// Move left back
	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyLeft})
	back := model.(PlayoffBracket)
	if back.cursorCol != 0 {
		t.Errorf("after Left from col 1, cursorCol = %d, want 0", back.cursorCol)
	}
}
