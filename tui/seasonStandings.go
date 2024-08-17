package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	//"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/app/datamodels"
	"log"
)

type seasonStandings struct {
	quitting  bool
	westTeams table.Model
	eastTeams table.Model
	height    int
	width     int
	maxHeight int
	maxWidth  int
	focused   bool
}

func (m seasonStandings) Init() tea.Cmd { return nil }

func initSeasonStandings(i list.Item, p *tea.Program) (*seasonStandings, error) {
	teams, headers, err := datamodels.PopulateTeamStats()
	if err != nil {
		log.Println("Could not populate player stats, error:", err)
		return nil, err
	}

	eastTeams, westTeams := teams.SplitStandingsPerConference()

	eastTeamsStrings := datamodels.ConvertToString(eastTeams)
	westTeamsStrings := datamodels.ConvertToString(westTeams)

	var (
		columns  []table.Column
		column   table.Column
		rows     []table.Row
		row      table.Row
		westRows []table.Row
		westRow  table.Row
	)

	for _, h := range headers {
		column = table.NewColumn(h, h, 15)
		columns = append(columns, column)
	}

	for _, r := range eastTeamsStrings {
		rowData := make(table.RowData)
		for i, rd := range r {
			columnTitle := columns[i].Title()
			rowData[columnTitle] = rd
		}
		row = table.NewRow(rowData)
		rows = append(rows, row)
	}

	for _, r := range westTeamsStrings {
		rowData := make(table.RowData)
		for i, rd := range r {
			columnTitle := columns[i].Title()
			rowData[columnTitle] = rd
		}
		westRow = table.NewRow(rowData)
		westRows = append(westRows, westRow)
	}

	eastTable := table.New(columns).WithRows(rows).WithMaxTotalWidth(120).Focused(true)

	westTable := table.New(columns).WithRows(westRows).WithMaxTotalWidth(120)

	m := &seasonStandings{eastTeams: eastTable, westTeams: westTable}

	return m, nil
}

func (m seasonStandings) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
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
				m.westTeams = m.westTeams.Focused(false)
				m.focused = !m.focused
			} else {
				m.eastTeams = m.eastTeams.Focused(false)
				m.westTeams = m.westTeams.Focused(true)
				m.focused = !m.focused
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

func (m seasonStandings) helpView() string {

	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m seasonStandings) View() string {
	if m.quitting {
		return ""
	}
	renderedEastTable := TableStyle.Render(m.eastTeams.View()) + "\n"
	renderedWestTable := TableStyle.Render(m.westTeams.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedEastTable, renderedWestTable, m.helpView())
	return DocStyle.Render(comboView)
}
