package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/datamodels"
	"log"
	"strings"
)

type LeagueLeaders struct {
	leaderboard table.Model
	quitting    bool
	height      int
	width       int
	maxHeight   int
	maxWidth    int
}

type fetchLeagueLeadersMsg struct {
	err     error
	columns []table.Column
	stats   []table.Row
}

func NewLeagueLeaders(size tea.WindowSizeMsg) (*LeagueLeaders, tea.Cmd, error) {
	m := &LeagueLeaders{
		height: size.Height,
		width:  size.Width,
	}

	_, _, err := datamodels.PopulatePlayerStats(datamodels.UnmarshallResponseJSON)
	if err != nil {
		return &LeagueLeaders{}, nil, fmt.Errorf("failed to populate player stats: %w", err)
	}

	cmd := fetchLeagueLeadersCmd()

	return m, cmd, nil
}

// fetchLeagueLeadersCmd fetches the required data from pre-existing JSON files and creates the table structure and rows
func fetchLeagueLeadersCmd() tea.Cmd {
	return func() tea.Msg {
		playerStats, headers, err := datamodels.PopulatePlayerStats(datamodels.UnmarshallResponseJSON)
		if err != nil {
			return fetchLeagueLeadersMsg{err: err}
		}

		playerStatsString := datamodels.ConvertToString(playerStats)

		var (
			columns []table.Column
			column  table.Column
			rows    []table.Row
			row     table.Row
		)

		for _, h := range headers {
			if !strings.Contains(h, "ID") {
				column = table.NewColumn(h, h, 15)
				columns = append(columns, column)
			}
		}

		for _, r := range playerStatsString {
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
			row = table.NewRow(rowData)
			rows = append(rows, row)
		}
		return fetchLeagueLeadersMsg{columns: columns, stats: rows}
	}
}

func (m LeagueLeaders) Init() tea.Cmd { return nil }

func (m LeagueLeaders) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case fetchLeagueLeadersMsg:
		if msg.err != nil {
			log.Println("fetch league leaders:", msg.err)
			return m, nil
		}

		t := table.New(msg.columns).
			WithRows(msg.stats).
			SelectableRows(true).
			WithHeaderVisibility(true).
			Focused(true).
			WithPageSize(20).
			WithMaxTotalWidth(145).
			WithHorizontalFreezeColumnCount(3).
			WithBaseStyle(lipgloss.NewStyle())

		m := &LeagueLeaders{leaderboard: t, maxHeight: 25, maxWidth: 125}
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, Keymap.Enter):
			selectedRows := m.leaderboard.SelectedRows()
			if len(selectedRows) == 1 {
				playerID := selectedRows[0].Data["PLAYER_ID"].(string)
				log.Println(playerID)
				//TODO: add player profile init logic
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

	m.leaderboard, cmd = m.leaderboard.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m LeagueLeaders) helpView() string {

	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m LeagueLeaders) View() string {
	if m.quitting {
		return ""
	}
	renderedLeaders := TableStyle.Render(m.leaderboard.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedLeaders, m.helpView())
	return DocStyle.Render(comboView)
}
