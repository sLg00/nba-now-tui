package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/datamodels"
	"log"
	"os"
	"reflect"
)

type boxScore struct {
	homeTeamBoxScore table.Model
	awayTeamBoxScore table.Model
	activeTable      int
	quitting         bool
	focused          bool
	width            int
	height           int
	maxWidth         int
	maxHeight        int
}

// getColsAndValues is a function to extract and filter fields and values from complex structs.
// Used currently on the boxScore rendering logic, to filter out unnecessary fields, YET
// leave the original data structures intact to preserve integrity and have the ability to extend and alter
func getColsAndValues(v interface{}) ([]string, []string) {
	var keys []string
	var values []string
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if field.Name == "NameI" || field.Name == "Position" || field.Name == "PersonId" {
			keys = append(keys, field.Name)
			if fieldValue.Kind() == reflect.String {
				values = append(values, fieldValue.String())
			} else {
				values = append(values, fmt.Sprintf("%v", fieldValue.Interface()))
			}

		} else if field.Name == "Statistics" {
			nestedKeys, nestedValues := extractStatistics(fieldValue.Interface())
			for _, nestedKey := range nestedKeys {
				keys = append(keys, fmt.Sprintf("%s", nestedKey))
			}
			values = append(values, nestedValues...)
		}
	}
	return keys, values
}

// extractStatistics is a companion function to getColsAndValues. It is used to handle extracting data
// from nested structures. It can be altered to be more generic in the future
func extractStatistics(v interface{}) ([]string, []string) {
	var keys []string
	var values []string

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Iterate through the fields of the Statistics struct
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		tag := field.Tag.Get("percentage")

		// Extract keys and values
		keys = append(keys, field.Name)
		if fieldValue.Kind() == reflect.Float64 {
			if tag == "true" {
				values = append(values, datamodels.FloatToPercent(fieldValue.Float()))
			}
		} else {
			values = append(values, fmt.Sprintf("%v", fieldValue.Interface()))
		}

	}

	return keys, values
}

func (m boxScore) Init() tea.Cmd { return nil }

// initBoxScore is the main function to render a game's box score, in very basic fashion.
// It takes a gameID as an input and uses it to access the corresponding game's box score,
// which is already present on the filesystem. The JSON files get downloaded when the Daily View is loaded.
func initBoxScore(gameID string, p *tea.Program) (*boxScore, error) {
	boxScoreData, err := datamodels.PopulateBoxScore(gameID)
	if err != nil {
		return nil, err
	}

	homeDataSet := boxScoreData.HomeTeam.BoxScorePlayers
	awayDataSet := boxScoreData.AwayTeam.BoxScorePlayers

	var column table.Column
	var columns []table.Column
	var homeRow table.Row
	var awayRow table.Row
	var homeRows []table.Row
	var awayRows []table.Row

	cols, _ := getColsAndValues(homeDataSet[0])

	columns = make([]table.Column, len(cols))
	for i, col := range cols {
		column = table.NewColumn(col, col, 15)
		columns[i] = column
	}

	for _, player := range homeDataSet {
		rowData := make(table.RowData)
		_, values := getColsAndValues(player)
		for i, value := range values {
			columnTitle := columns[i].Title()
			rowData[columnTitle] = value
		}
		homeRow = table.NewRow(rowData)
		homeRows = append(homeRows, homeRow)
	}

	for _, player := range awayDataSet {
		rowData := make(table.RowData)
		_, values := getColsAndValues(player)
		for i, value := range values {
			columnTitle := columns[i].Title()
			rowData[columnTitle] = value
		}
		awayRow = table.NewRow(rowData)
		awayRows = append(awayRows, awayRow)
	}

	homeTable := table.New(columns).
		WithRows(homeRows).
		SelectableRows(true).
		WithMaxTotalWidth(120).
		Focused(true)

	awayTable := table.New(columns).
		WithRows(awayRows).
		SelectableRows(true).
		WithMaxTotalWidth(120).
		Focused(false)

	m := &boxScore{
		homeTeamBoxScore: homeTable,
		awayTeamBoxScore: awayTable}

	return m, nil
}

func (m boxScore) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	var selectedRows []table.Row
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			dv, err := initDailyView()
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			return dv.Update(WindowSize)
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, Keymap.Tab):
			if m.focused {
				m.homeTeamBoxScore = m.homeTeamBoxScore.Focused(true)
				m.activeTable = 0
				m.awayTeamBoxScore = m.awayTeamBoxScore.Focused(false)
			} else {
				m.homeTeamBoxScore = m.homeTeamBoxScore.Focused(false)
				m.activeTable = 1
				m.awayTeamBoxScore = m.awayTeamBoxScore.Focused(true)
			}
			m.focused = !m.focused
		case key.Matches(msg, Keymap.Enter):
			if m.activeTable == 0 {
				selectedRows = m.homeTeamBoxScore.SelectedRows()
			} else {
				selectedRows = m.awayTeamBoxScore.SelectedRows()
			}
			if len(selectedRows) == 1 {
				personId := selectedRows[0].Data["PersonId"].(string)
				log.Println(personId)
				//TODO: add player profile view init
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
	m.homeTeamBoxScore, cmd = m.homeTeamBoxScore.Update(msg)
	m.awayTeamBoxScore, cmd = m.awayTeamBoxScore.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m boxScore) helpView() string {
	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m boxScore) View() string {
	if m.quitting {
		return ""
	}
	renderedHomeBoxScore := TableStyle.Render(m.homeTeamBoxScore.View()) + "\n"
	renderedAwayBoxScore := TableStyle.Render(m.awayTeamBoxScore.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedHomeBoxScore, renderedAwayBoxScore, m.helpView())
	return DocStyle.Render(comboView)
}
