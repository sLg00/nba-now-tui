package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sLg00/nba-now-tui/app/datamodels"
	"os"
)

func createLeagueLeadersTable() table.Model {
	playerStats, headers, err := datamodels.PopulatePlayerStats()
	if err != nil {
		fmt.Println("error:", err)
	}
	playerStatsString := playerStats.ConvertToString()

	var (
		columns []table.Column
		column  table.Column
		rows    []table.Row
		row     table.Row
	)

	for _, h := range headers {
		column = table.Column{
			Title: h,
			Width: 10,
		}
		columns = append(columns, column)
	}

	for _, r := range playerStatsString {
		row = r
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(25),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func RunTUI() {
	t := createLeagueLeadersTable()
	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
