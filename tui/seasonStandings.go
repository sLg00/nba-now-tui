package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/converters"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
)

// SeasonStandings is the main model of the season standings view, containing all the relevant fields
type SeasonStandings struct {
	quitting    bool
	westTeams   table.Model
	eastTeams   table.Model
	activeTable int
	height      int
	width       int
	maxHeight   int
	maxWidth    int
	focused     bool
}

// fetchSeasonStandingsMsg is a structure to transform the raw data into the season standings tables
type fetchSeasonStandingsMsg struct {
	err       error
	eastTable table.Model
	westTable table.Model
}

// teamProfileDownloadedMsg is returned by te downloadProfile function to note whether the API call
// to retrieve team details has succeeded and pass on the teamID attribute for futher processing
type teamProfileDownloadedMsg struct {
	err    error
	teamID string
}

// downloadProfile executes the call to MakeOnDemandRequests using teamID as an input parameter.
// It returns a teamProfileDownloadedMsg command
func downloadProfile(teamID string) tea.Cmd {
	return func() tea.Msg {
		err := nbaAPI.NewClient().FetchTeamProfile(teamID)
		return teamProfileDownloadedMsg{err: err, teamID: teamID}
	}
}

func NewSeasonStandings(size tea.WindowSizeMsg) (*SeasonStandings, tea.Cmd, error) {
	m := &SeasonStandings{
		height: size.Height,
		width:  size.Width,
	}
	cmd := fetchSeasonStandingsCmd()

	return m, cmd, nil
}

// fetchSeasonStandingsCmd fetches and prepares the data required to display the season standings tables
func fetchSeasonStandingsCmd() tea.Cmd {
	return func() tea.Msg {
		cl, err := nbaAPI.NewClient().Loader.LoadSeasonStandings()
		if err != nil {
			log.Println("Error loading season standings:", err)
		}
		teams, headers, err := converters.PopulateTeamStats(cl)
		if err != nil {
			return fetchSeasonStandingsMsg{err: err}
		}

		eastTeams, westTeams := teams.SplitStandingsPerConference()

		eastTeamsStrings := types.ConvertToStringMatrix(eastTeams)
		westTeamsStrings := types.ConvertToStringMatrix(westTeams)

		eastTable := buildTables(headers, eastTeamsStrings, types.Team{}).
			SelectableRows(true).
			Focused(true)
		westTable := buildTables(headers, westTeamsStrings, types.Team{}).
			SelectableRows(true)

		return fetchSeasonStandingsMsg{
			eastTable: eastTable,
			westTable: westTable}
	}
}

func (m SeasonStandings) Init() tea.Cmd { return nil }

func (m SeasonStandings) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	var selectedRows []table.Row
	switch msg := msg.(type) {
	case fetchSeasonStandingsMsg:
		if msg.err != nil {
			log.Println("could not fetch season standings:", msg.err)
			return m, nil
		}
		eastTable := msg.eastTable

		westTable := msg.westTable
		m := &SeasonStandings{eastTeams: eastTable, westTeams: westTable}
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, Keymap.Tab):
			if m.focused {
				m.eastTeams = m.eastTeams.Focused(true)
				m.activeTable = 0
				m.westTeams = m.westTeams.Focused(false)
				m.focused = !m.focused
			} else {
				m.eastTeams = m.eastTeams.Focused(false)
				m.activeTable = 1
				m.westTeams = m.westTeams.Focused(true)
				m.focused = !m.focused
			}
		case key.Matches(msg, Keymap.Enter):
			if m.activeTable == 0 {
				selectedRows = m.eastTeams.SelectedRows()
			} else {
				selectedRows = m.westTeams.SelectedRows()
			}
			if len(selectedRows) == 1 {
				teamID := selectedRows[0].Data["TeamID"].(string)
				log.Println(teamID)
				dlp := downloadProfile(teamID)
				return m, dlp
			}
			if len(selectedRows) > 1 || len(selectedRows) < 1 {
				log.Println("Either 0 rows or more than 1 row were selected")
				//TODO: Display pop-up with User error! :)
			}
		}
	case teamProfileDownloadedMsg:
		if msg.err != nil {
			log.Println("could not download team profile:", msg.err)
		}
		tp, cmd, err := NewTeamProfile(msg.teamID, WindowSize)
		if err != nil {
			log.Println("could not load team profile:", err)
			//TODO: add error modal
			return InitMenu()
		}
		return tp, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.width > m.maxWidth {
			m.maxWidth = m.width
		}
		if m.height > m.maxHeight {
			m.maxHeight = m.height
		}
	}
	m.eastTeams, cmd = m.eastTeams.Update(msg)
	m.westTeams, cmd = m.westTeams.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m SeasonStandings) helpView() string {

	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m SeasonStandings) View() string {
	if m.quitting {
		return ""
	}
	renderedEastTable := m.eastTeams.View() + "\n"
	renderedWestTable := m.westTeams.View() + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n",
		"<< EAST >> \n",
		renderedEastTable,
		"<< WEST >> \n",
		renderedWestTable,
		m.helpView())
	return DocStyle.Render(comboView)
}
