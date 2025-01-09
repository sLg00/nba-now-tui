package tui

import (
	"bytes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/sLg00/nba-now-tui/cmd/helpers"
	"testing"
)

func TestMenuItemCreation(t *testing.T) {
	items, err := createMenuItems()
	if err != nil {
		t.Fatalf("createMenuItems() returned unexpected error: %v", err)
	}

	if len(items) != 5 {
		t.Errorf("createMenuItems() expected 5 menu items , got  %d", len(items))
	}

	expectedTitles := []string{
		"Daily Scores",
		"Season Standings",
		"League Leaders",
		"[N/A] Recent News",
		"[N/A] Playoff Bracket",
	}

	for i, item := range items {
		menuItem := item.(menuItem)
		if menuItem.Index() != i {
			t.Errorf("createMenuItems() expected index %d , got  %d", i, menuItem.Index())
		}
		if menuItem.Title() != expectedTitles[i] {
			t.Errorf("Expected title %s, got %s", expectedTitles[i], menuItem.Title())
		}
	}
}

func TestModelStateManagement(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	model, _ := InitMenu()
	m := model.(Model)

	if m.quitting || m.requestsMade {
		t.Error("new model should be instantiated with false flags")
	}

	cmd := m.Init()
	if cmd == nil {
		t.Error("init command should not be nil on first load")
	}

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	updatedModel := newModel.(Model)
	if !updatedModel.quitting {
		t.Error("model should be in quitting state after Q")
	}

}

func TestMenuItemSelection(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	model, _ := InitMenu()
	m := model.(Model)

	m.menu.Select(2)
	selected := m.menu.SelectedItem()

	if selected.FilterValue() != "League Leaders" {
		t.Errorf("expected League Leaders got %s", selected.FilterValue())
	}

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("cmd should not be nil after selecting menu item")
	}

	if _, ok := newModel.(Model); ok {
		t.Error("expected different model type after selection")
	}
}

func TestMenuDisplay(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	model, _ := InitMenu()
	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, testModel.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("NBA on"))
	})

	testModel.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	testModel.WaitFinished(t)
}
