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

type TeamProfile struct {
	teamTable table.Model
	width     int
	height    int
}

type teamProfileFetchedMsg struct {
	err         error
	teamProfile table.Model
}

func NewTeamProfile(teamID string, size tea.WindowSizeMsg) (*TeamProfile, tea.Cmd, error) {
	m := &TeamProfile{
		width:  size.Width,
		height: size.Height,
	}

	cl, err := nbaAPI.NewNewClient().Loader.LoadTeamProfile(teamID)
	_, _, err = converters.PopulateTeamInfo(cl)
	if err != nil {
		return &TeamProfile{}, nil, err
	}

	cmd := fetchTeamProfileMsg(teamID)

	return m, cmd, nil
}

func fetchTeamProfileMsg(teamID string) tea.Cmd {
	return func() tea.Msg {
		cl, err := nbaAPI.NewNewClient().Loader.LoadTeamProfile(teamID)
		if err != nil {
			log.Println("error loading team profile:", err)
		}
		data, headers, err := converters.PopulateTeamInfo(cl)
		if err != nil {
			return teamProfileFetchedMsg{err: err}
		}

		teamProfileStrings := types.ConvertToStringFlat(data)

		var column table.Column
		var columns []table.Column
		var row table.Row
		var rows []table.Row

		for _, h := range headers {
			column = table.NewColumn(h, h, 20)
			columns = append(columns, column)
		}

		usedRow := teamProfileStrings
		rowData := make(table.RowData)
		for i, cell := range usedRow {
			headerName := headers[i]
			rowData[headerName] = cell
			row = table.NewRow(rowData)
		}
		rows = append(rows, row)

		teamTable := table.New(columns).WithRows(rows)

		return teamProfileFetchedMsg{err: nil, teamProfile: teamTable}
	}
}

func (m TeamProfile) Init() tea.Cmd { return nil }

func (m TeamProfile) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case teamProfileDownloadedMsg:
		if msg.err != nil {
			log.Println("could not download team profile:", msg.err)
			return m, nil
		}
	case teamProfileFetchedMsg:
		if msg.err != nil {
			log.Println("could not load team profile:", msg.err)
			return m, nil
		}
		m := &TeamProfile{teamTable: msg.teamProfile}
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.teamTable, cmd = m.teamTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m TeamProfile) helpView() string {

	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m TeamProfile) View() string {
	renderedTable := TableStyle.Render(m.teamTable.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedTable, m.helpView())
	return DocStyle.Render(comboView)
}
