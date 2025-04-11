package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/converters"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type InstantiatedBoxScore struct {
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

type boxScoreFetchedMsg struct {
	err                  error
	boxScoreTableColumns []table.Column
	homeBoxScoreData     []table.Row
	awayBoxScoreData     []table.Row
}

// NewBoxScore is a factory function to instantiate a BoxScore when the BoxScore is opened from the Daily View.
// It takes gameId and WindowSize as inputs and returns a model, command and error
func NewBoxScore(gameId string, size tea.WindowSizeMsg) (*InstantiatedBoxScore, tea.Cmd, error) {
	m := &InstantiatedBoxScore{
		width:  size.Width,
		height: size.Height,
	}

	cl, err := nbaAPI.NewClient().Loader.LoadBoxScore(gameId)
	_, err = converters.PopulateBoxScore(cl)
	if err != nil {
		return &InstantiatedBoxScore{}, nil, fmt.Errorf("failed to populate box score: %w", err)
	}

	cmd := fetchBoxSoresCmd(gameId)

	return m, cmd, nil
}

// fetchBoxScoresCmd "fetches" and processes the given game data to eventually render a box score
func fetchBoxSoresCmd(gameID string) tea.Cmd {
	return func() tea.Msg {
		cl, err := nbaAPI.NewClient().Loader.LoadBoxScore(gameID)
		if err != nil {
			log.Printf("failed to load box score for game id %s: %v", gameID, err)
		}
		boxScoreData, err := converters.PopulateBoxScore(cl)
		if err != nil {
			return boxScoreFetchedMsg{err: err}
		}
		homeDataSet := boxScoreData.HomeTeam.BoxScorePlayers
		awayDataSet := boxScoreData.AwayTeam.BoxScorePlayers

		var column table.Column
		var columns []table.Column
		var homeRow table.Row
		var awayRow table.Row
		var homeRows []table.Row
		var awayRows []table.Row

		cols, _, widths := getColsAndValues(homeDataSet[0])

		columns = make([]table.Column, len(cols))
		for i, col := range cols {
			width := 15
			if i < len(widths) {
				width = widths[i]
			}
			column = table.NewColumn(col, col, width)
			columns[i] = column
		}

		for _, player := range homeDataSet {
			rowData := make(table.RowData)
			_, values, _ := getColsAndValues(player)
			for i, value := range values {
				columnTitle := columns[i].Title()
				rowData[columnTitle] = value
			}
			homeRow = table.NewRow(rowData)
			homeRows = append(homeRows, homeRow)
		}

		for _, player := range awayDataSet {
			rowData := make(table.RowData)
			_, values, _ := getColsAndValues(player)
			for i, value := range values {
				columnTitle := columns[i].Title()
				rowData[columnTitle] = value
			}
			awayRow = table.NewRow(rowData)
			awayRows = append(awayRows, awayRow)
		}

		return boxScoreFetchedMsg{
			boxScoreTableColumns: columns,
			homeBoxScoreData:     homeRows,
			awayBoxScoreData:     awayRows,
		}
	}
}

// getColsAndValues is a function to extract and filter fields and values from complex structs.
// Used currently on the boxScore rendering logic, to filter out unnecessary fields, YET
// leave the original data structures intact to preserve integrity and have the ability to extend and alter
func getColsAndValues(v interface{}) ([]string, []string, []int) {
	var keys []string
	var values []string
	var widths []int

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		isVisibleTag := field.Tag.Get("isVisible")
		if isVisibleTag == "false" {
			continue
		}

		if field.Name == "NameI" || field.Name == "Position" || field.Name == "PersonId" {

			displayName := field.Tag.Get("display")
			if displayName == "" {
				displayName = field.Name
			}
			keys = append(keys, displayName)

			width := 15
			widthTag := field.Tag.Get("width")
			if widthTag != "" {
				if w, err := strconv.Atoi(widthTag); err == nil {
					width = w
				}
			}
			widths = append(widths, width)

			if fieldValue.Kind() == reflect.String {
				values = append(values, fieldValue.String())
			} else {
				values = append(values, fmt.Sprintf("%v", fieldValue.Interface()))
			}

		} else if field.Name == "Statistics" {
			nestedKeys, nestedValues, nestedWidths := extractStatistics(fieldValue.Interface())
			for i, nestedKey := range nestedKeys {
				keys = append(keys, fmt.Sprintf("%s", nestedKey))
				widths = append(widths, nestedWidths[i])
			}
			values = append(values, nestedValues...)
		}
	}
	return keys, values, widths
}

// extractStatistics is a companion function to getColsAndValues. It is used to handle extracting data
// from nested structures. It can be altered to be more generic in the future
func extractStatistics(v interface{}) ([]string, []string, []int) {
	var keys []string
	var values []string
	var widths []int

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Iterate through the fields of the Statistics struct
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Check visibility
		isVisibleTag := field.Tag.Get("isVisible")
		if isVisibleTag == "false" {
			continue // Skip invisible fields
		}

		// Get display name
		displayName := field.Tag.Get("display")
		jsonTag := field.Tag.Get("json")
		if displayName == "" && jsonTag != "" {
			tagParts := strings.Split(jsonTag, ",")
			displayName = tagParts[0]
		}
		if displayName == "" {
			displayName = field.Name
		}

		// Get column width
		width := 15 // Default width
		widthTag := field.Tag.Get("width")
		if widthTag != "" {
			if w, err := strconv.Atoi(widthTag); err == nil {
				width = w
			}
		}

		// Add key, width, and format value
		keys = append(keys, displayName)
		widths = append(widths, width)

		// Format value (handle percentages, etc.)
		tag := field.Tag.Get("percentage")
		if fieldValue.Kind() == reflect.Float64 && tag == "true" {
			values = append(values, types.FloatToPercent(fieldValue.Float()))
		} else {
			values = append(values, fmt.Sprintf("%v", fieldValue.Interface()))
		}
	}

	return keys, values, widths
}

func (m InstantiatedBoxScore) Init() tea.Cmd { return nil }

func (m InstantiatedBoxScore) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	var selectedRows []table.Row
	switch msg := msg.(type) {
	case boxScoreFetchedMsg:
		if msg.err != nil {
			log.Println("error fetching box score:", msg.err)
			return m, nil
		}
		homeTable := table.New(msg.boxScoreTableColumns).
			WithRows(msg.homeBoxScoreData).
			SelectableRows(true).
			WithMaxTotalWidth(140).
			Focused(true).
			WithHorizontalFreezeColumnCount(3)

		awayTable := table.New(msg.boxScoreTableColumns).
			WithRows(msg.awayBoxScoreData).
			SelectableRows(true).
			WithMaxTotalWidth(140).
			Focused(false).
			WithHorizontalFreezeColumnCount(3)

		m := &InstantiatedBoxScore{
			homeTeamBoxScore: homeTable,
			awayTeamBoxScore: awayTable}
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			dv, cmd, err := NewDailyView(WindowSize)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			return dv, cmd
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

func (m InstantiatedBoxScore) helpView() string {
	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m InstantiatedBoxScore) View() string {
	if m.quitting {
		return ""
	}
	renderedHomeBoxScore := TableStyle.Render(m.homeTeamBoxScore.View()) + "\n"
	renderedAwayBoxScore := TableStyle.Render(m.awayTeamBoxScore.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedHomeBoxScore, renderedAwayBoxScore, m.helpView())
	return DocStyle.Render(comboView)
}
