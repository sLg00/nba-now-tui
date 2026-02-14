package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"log"
	"os"
)

// Model acts as the main model of the TUI. It's just to build the initial menu
type Model struct {
	menu         list.Model
	quitting     bool
	requestsMade bool
}

// menuItem is a singular list item within a list
type menuItem struct {
	index              int
	title, description string
}

// requestsFinishedMsg is a struct for a custom tea.Msg that denotes if the default requests and subsequent
// storing of JSON was successful
type requestsFinishedMsg struct {
	err error
}

// these methods on the menuItem model ensure that the menuItem objects satisfies the requirements of list.Model struct
func (m menuItem) Title() string       { return m.title }
func (m menuItem) Index() int          { return m.index }
func (m menuItem) Description() string { return m.description }
func (m menuItem) FilterValue() string { return m.title }

// createMenuItems returns a []list.Item with 1-N items to be displayed in a list
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

// makeDefaultRequests is used to initiate default api requests and results storing on TUI launch.
// It's only ran once when the app starts. Subsequent returns to the main menu do not trigger it again.
func makeInitialRequests() tea.Cmd {
	return func() tea.Msg {
		err := nbaAPI.NewClient().MakeDefaultRequests()
		return requestsFinishedMsg{err: err}
	}
}

// InitMenu creates the list object and returns the model
func InitMenu() (tea.Model, tea.Cmd) {
	items, err := createMenuItems()
	m := Model{
		menu:         list.New(items, list.NewDefaultDelegate(), 10, 10),
		requestsMade: false,
	}
	if WindowSize.Height != 0 {
		top, right, bottom, left := DocStyle.GetMargin()
		m.menu.SetSize(WindowSize.Width-left-right, WindowSize.Height-top-bottom-1)
	}
	currentDate, _ := nbaAPI.NewClient().Dates.GetCurrentDate()
	m.menu.Title = "NBA on " + currentDate
	m.menu.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			Keymap.Back,
			Keymap.Enter,
			Keymap.Tab,
		}
	}
	return m, func() tea.Msg { return errMsg{err} }
}

// Init is in charge of executing the initialization logic,
// in the case of mainMenu, it fires off the default requests
func (m Model) Init() tea.Cmd {
	if !m.requestsMade {
		m.requestsMade = true
		return tea.Batch(makeInitialRequests())
	}
	return nil
}

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
			selectedItem := m.menu.SelectedItem()
			switch {
			case selectedItem.FilterValue() == "League Leaders":
				ll, cmd, err := NewLeagueLeaders(WindowSize)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				return ll, cmd
			case selectedItem.FilterValue() == "Daily Scores":
				dv, cmd, err := NewDailyView(WindowSize)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				return dv, cmd
			case selectedItem.FilterValue() == "Season Standings":
				ss, cmd, err := NewSeasonStandings(WindowSize)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				return ss, cmd
			case selectedItem.FilterValue() == "Recent News":
				nv, cmd, err := NewNewsView(WindowSize)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				return nv, cmd
			}
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		default:
			m.menu, cmd = m.menu.Update(msg)
		}
		cmds = append(cmds, cmd)

	case requestsFinishedMsg:
		if msg.err != nil {
			log.Println("Baseline data population finished with error", msg.err)
		} else {
			log.Println("Baseline data population finished successfully")
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return DocStyle.Render(m.menu.View() + "\n")
}
