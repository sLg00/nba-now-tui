package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	menu     list.Model
	styles   Styles
	quitting bool
}

type menuItem struct {
	index              int
	title, description string
}

func (m menuItem) Title() string       { return m.title }
func (m menuItem) Index() int          { return m.index }
func (m menuItem) Description() string { return m.description }
func (m menuItem) FilterValue() string { return m.title }

func createMenuItems() ([]list.Item, error) {
	items := []list.Item{
		menuItem{
			index:       0,
			title:       "Daily Scores",
			description: "All game results for a given date",
		}, menuItem{
			index:       1,
			title:       "Season Standings",
			description: "Regular season standings",
		}, menuItem{
			index:       2,
			title:       "League Leaders",
			description: "All players sorted by PPG",
		}, menuItem{
			index:       3,
			title:       "Recent News",
			description: "Headlines from around the league",
		}}
	return items, nil
}

func InitMenu() (tea.Model, tea.Cmd) {
	items, err := createMenuItems()
	m := Model{
		menu:   list.New(items, list.NewDefaultDelegate(), 10, 10),
		styles: ListStyles(),
	}
	if WindowSize.Height != 0 {
		top, right, bottom, left := DocStyle.GetMargin()
		m.menu.SetSize(WindowSize.Width-left-right, WindowSize.Height-top-bottom-1)
	}

	m.menu.Title = "NBA on " + Date()
	m.menu.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			Keymap.Back,
			Keymap.Quit,
			Keymap.Enter,
		}
	}
	return m, func() tea.Msg { return errMsg{err} }
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		WindowSize = msg
		top, right, bottom, left := DocStyle.GetMargin()
		m.menu.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Enter):
			//selectedItem := m.menu.SelectedItem()
			//subPage := InitSubPage(selectedItem, Program)
			//return subPage.Update(WindowSize)
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		default:
			m.menu, cmd = m.menu.Update(msg)
		}
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return DocStyle.Render(m.menu.View() + "\n")
}
