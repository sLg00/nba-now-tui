package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/sLg00/nba-now-tui/app/datamodels"
	"slices"
	"strings"
)

// removeIndex is a helper function to delete elements from a slice
func removeIndex[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

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
		if strings.Contains(column.Title, "PLAYER") {
			column.Width = 25
		}
		columns = append(columns, column)
	}

	for _, r := range playerStatsString {
		row = r
		//row = removeIndex(row, 0)
		//row = removeIndex(row, 2)
		//delete ID slices
		row = slices.Delete(row, 0, 1)
		row = slices.Delete(row, 2, 3)
		rows = append(rows, row)
	}

	// or loop to filter out ID column headers
	var filteredColumns []table.Column
	for _, c := range columns {
		if !strings.Contains(c.Title, "ID") {
			filteredColumns = append(filteredColumns, c)
		}
	}

	t := table.New(
		table.WithColumns(filteredColumns),
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
