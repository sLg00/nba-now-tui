package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sLg00/nba-now-tui/app/datamodels"
	"log"
	"slices"
	"strings"
)

// removeIndex is a helper function to delete elements from a slice
func removeIndex[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

type leagueLeaders struct {
	leaderboard table.Model
	quitting    bool
}

func (m leagueLeaders) Init() tea.Cmd { return nil }

// initLeagueLeadersTable returns a table.Model which is populated with the current league leaders (per PPG)
func initLeagueLeaders(i list.Item, p *tea.Program) (*leagueLeaders, error) {
	playerStats, headers, err := datamodels.PopulatePlayerStats()
	if err != nil {
		log.Println("Could not populate player stats, error:", err)
		return nil, err
	}
	playerStatsString := datamodels.ConvertToString(playerStats)

	var (
		columns    []table.Column
		column     table.Column
		rows       []table.Row
		row        table.Row
		sortedRows []table.Row
	)

	for _, h := range headers {
		column = table.Column{
			Title: h,
			Width: 10,
		}
		if strings.Contains(column.Title, "PLAYER") {
			column.Width = 25
		}
		columns = append(columns, column)
	}

	// filter out "ID" slices
	for _, r := range playerStatsString {
		row = r
		//row = removeIndex(row, 0)
		//row = removeIndex(row, 2)
		//delete ID slices
		row = slices.Delete(row, 0, 1)
		row = slices.Delete(row, 2, 3)
		rows = append(rows, row)
	}

	//sort the data in a logical order with major categories at the front
	for _, r := range rows {
		row = r
		var sortedRow table.Row
		sortedRow = append(sortedRow, row[0:5]...)
		sortedRow = append(sortedRow, row[21:23]...)
		sortedRow = append(sortedRow, row[17:21]...)
		sortedRow = append(sortedRow, row[5:17]...)
		sortedRows = append(sortedRows, sortedRow)
	}

	// for loop to filter out ID column headers
	var filteredColumns []table.Column
	for _, c := range columns {
		if !strings.Contains(c.Title, "ID") {
			filteredColumns = append(filteredColumns, c)
		}
	}

	//sort columns in a more logical order, with major categories at the beginning
	var sortedColumns []table.Column
	sortedColumns = append(sortedColumns, filteredColumns[0:5]...)
	sortedColumns = append(sortedColumns, filteredColumns[21:23]...)
	sortedColumns = append(sortedColumns, filteredColumns[17:21]...)
	sortedColumns = append(sortedColumns, filteredColumns[5:17]...)

	t := table.New(
		table.WithColumns(sortedColumns),
		table.WithRows(sortedRows),
		table.WithFocused(true),
		table.WithHeight(25),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	m := &leagueLeaders{leaderboard: t}

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
		}
	}
	m.leaderboard, cmd = m.leaderboard.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m leagueLeaders) helpView() string {
	// TODO: use the keymaps to populate the help string
	return HelpStyle("\n ↑/↓: navigate  • backspace: back • q: quit\n")
}

func (m leagueLeaders) View() string {
	if m.quitting {
		return ""
	}
	renderedLeaders := TableStyle.Render(m.leaderboard.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedLeaders, m.helpView())
	return DocStyle.Render(comboView)
}
