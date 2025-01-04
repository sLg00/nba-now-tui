package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"strings"

	//"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/datamodels"
	"log"
)

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

type fetchSeasonStandingsMsg struct {
	err           error
	columns       []table.Column
	eastTeamStats []table.Row
	westTeamStats []table.Row
}

func NewSeasonStandings(size tea.WindowSizeMsg) (*SeasonStandings, tea.Cmd, error) {
	m := &SeasonStandings{
		height: size.Height,
		width:  size.Width,
	}

	_, _, err := datamodels.PopulateTeamStats(datamodels.UnmarshallResponseJSON)
	if err != nil {
		return &SeasonStandings{}, nil, fmt.Errorf("could not populate team stats: %w", err)
	}

	cmd := fetchSeasonStandingsCmd()

	return m, cmd, nil
}

func fetchSeasonStandingsCmd() tea.Cmd {
	return func() tea.Msg {
		teams, headers, err := datamodels.PopulateTeamStats(datamodels.UnmarshallResponseJSON)
		if err != nil {
			return fetchSeasonStandingsMsg{err: err}
		}

		eastTeams, westTeams := teams.SplitStandingsPerConference()

		eastTeamsStrings := datamodels.ConvertToString(eastTeams)
		westTeamsStrings := datamodels.ConvertToString(westTeams)

		var (
			columns  []table.Column
			column   table.Column
			eastRows []table.Row
			eastRow  table.Row
			westRows []table.Row
			westRow  table.Row
		)

		for _, h := range headers {
			if !strings.Contains(h, "ID") {
				column = table.NewColumn(h, h, 15)
				columns = append(columns, column)
			}
		}

		for _, r := range eastTeamsStrings {
			rowData := make(table.RowData)
			visibleColumnIndex := 0
			for i, rd := range r {
				headerName := headers[i]
				if strings.Contains(headerName, "ID") {
					rowData[headerName] = rd
				} else {
					columnTitle := columns[visibleColumnIndex].Title()
					rowData[columnTitle] = rd
					visibleColumnIndex++
				}
			}
			eastRow = table.NewRow(rowData)
			eastRows = append(eastRows, eastRow)
		}

		for _, r := range westTeamsStrings {
			rowData := make(table.RowData)
			visibleColumnIndex := 0
			for i, rd := range r {
				headerName := headers[i]
				if strings.Contains(headerName, "ID") {
					rowData[headerName] = rd
				} else {
					columnTitle := columns[visibleColumnIndex].Title()
					rowData[columnTitle] = rd
					visibleColumnIndex++
				}
			}
			westRow = table.NewRow(rowData)
			westRows = append(westRows, westRow)
		}

		return fetchSeasonStandingsMsg{
			columns:       columns,
			eastTeamStats: eastRows,
			westTeamStats: westRows}
	}
}

func (m SeasonStandings) Init() tea.Cmd { return nil }

//// initSeasonStandings gets, filters and populates the Season Standings tables
//func initSeasonStandings(i list.Item, p *tea.Program) (*SeasonStandings, error) {
//	teams, headers, err := datamodels.PopulateTeamStats(datamodels.UnmarshallResponseJSON)
//	if err != nil {
//		log.Println("Could not populate player stats, error:", err)
//		return nil, err
//	}
//
//	eastTeams, westTeams := teams.SplitStandingsPerConference()
//
//	eastTeamsStrings := datamodels.ConvertToString(eastTeams)
//	westTeamsStrings := datamodels.ConvertToString(westTeams)
//
//	var (
//		columns  []table.Column
//		column   table.Column
//		rows     []table.Row
//		row      table.Row
//		westRows []table.Row
//		westRow  table.Row
//	)
//
//	for _, h := range headers {
//		if !strings.Contains(h, "ID") {
//			column = table.NewColumn(h, h, 15)
//			columns = append(columns, column)
//		}
//	}
//
//	for _, r := range eastTeamsStrings {
//		rowData := make(table.RowData)
//		visibleColumnIndex := 0
//		for i, rd := range r {
//			headerName := headers[i]
//			if strings.Contains(headerName, "ID") {
//				rowData[headerName] = rd
//			} else {
//				columnTitle := columns[visibleColumnIndex].Title()
//				rowData[columnTitle] = rd
//				visibleColumnIndex++
//			}
//		}
//		row = table.NewRow(rowData)
//		rows = append(rows, row)
//	}
//
//	for _, r := range westTeamsStrings {
//		rowData := make(table.RowData)
//		visibleColumnIndex := 0
//		for i, rd := range r {
//			headerName := headers[i]
//			if strings.Contains(headerName, "ID") {
//				rowData[headerName] = rd
//			} else {
//				columnTitle := columns[visibleColumnIndex].Title()
//				rowData[columnTitle] = rd
//				visibleColumnIndex++
//			}
//		}
//		westRow = table.NewRow(rowData)
//		westRows = append(westRows, westRow)
//	}
//
//	eastTable := table.New(columns).
//		WithRows(rows).
//		SelectableRows(true).
//		WithMaxTotalWidth(120).
//		Focused(true)
//
//	westTable := table.New(columns).
//		WithRows(westRows).
//		SelectableRows(true).
//		WithMaxTotalWidth(120)
//
//	m := &SeasonStandings{eastTeams: eastTable, westTeams: westTable}
//
//	return m, nil
//}

func (m SeasonStandings) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	var selectedRows []table.Row
	switch msg := msg.(type) {
	case fetchSeasonStandingsMsg:
		if msg.err != nil {
			log.Println("could not fetch season standings:", msg.err)
			return m, nil
		}
		eastTable := table.New(msg.columns).
			WithRows(msg.eastTeamStats).
			SelectableRows(true).
			WithMaxTotalWidth(120).
			Focused(true)

		westTable := table.New(msg.columns).
			WithRows(msg.westTeamStats).
			SelectableRows(true).
			WithMaxTotalWidth(120)

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
				//TODO: add team profile init logic
			}
			if len(selectedRows) > 1 || len(selectedRows) < 1 {
				log.Println("Either 0 rows or more than 1 row were selected")
				//TODO: Display pop-up with User error! :)
			}
		}
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
	renderedEastTable := TableStyle.Render(m.eastTeams.View()) + "\n"
	renderedWestTable := TableStyle.Render(m.westTeams.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedEastTable, renderedWestTable, m.helpView())
	return DocStyle.Render(comboView)
}
