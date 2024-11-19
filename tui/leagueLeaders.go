package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/datamodels"
	"log"
	"strings"
)

// removeIndex is a helper function to delete elements from a slice
func removeIndex[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

type leagueLeaders struct {
	leaderboard table.Model
	quitting    bool
	height      int
	width       int
	maxHeight   int
	maxWidth    int
}

func (m leagueLeaders) Init() tea.Cmd { return nil }

// initLeagueLeadersTable returns a table.Model which is populated with the current league leaders (PPG)
func initLeagueLeaders(i list.Item, p *tea.Program) (*leagueLeaders, error) {
	playerStats, headers, err := datamodels.PopulatePlayerStats()
	if err != nil {
		log.Println("Could not populate player stats, error:", err)
		return nil, err
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

	t := table.New(columns).
		WithRows(rows).
		SelectableRows(true).
		WithHeaderVisibility(true).
		Focused(true).
		WithPageSize(20).
		WithMaxTotalWidth(150).
		WithHorizontalFreezeColumnCount(2).
		WithBaseStyle(lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()))

	m := &leagueLeaders{leaderboard: t, maxHeight: 25, maxWidth: 80}

	return m, nil
}

func (m leagueLeaders) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
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

func (m leagueLeaders) helpView() string {

	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m leagueLeaders) View() string {
	if m.quitting {
		return ""
	}
	renderedLeaders := TableStyle.Render(m.leaderboard.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedLeaders, m.helpView())
	return DocStyle.Render(comboView)
}
