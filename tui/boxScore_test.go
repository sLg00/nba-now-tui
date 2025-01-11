package tui

import (
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/helpers"
	"reflect"
	"testing"
)

type mockBoxScorePlayer struct {
	NameI      string
	Position   string
	PersonId   string
	Statistics mockStats
}

type mockStats struct {
	Points   int
	Rebounds int
	Assists  int
	Fg3ptPct float64 `percentage:"true"`
}

type mockBoxScoreFetchedMsg struct {
	err                  error
	boxScoreTableColumns []table.Column
	homeBoxScoreData     []table.Row
	awayBoxScoreData     []table.Row
}

func TestGetColsAndValues(t *testing.T) {
	player := mockBoxScorePlayer{
		NameI:    "Rasheed Wallace",
		Position: "F",
		PersonId: "123456",
		Statistics: mockStats{
			Points:   33,
			Rebounds: 18,
			Assists:  7,
			Fg3ptPct: 0.388,
		},
	}

	expectedCols := []string{
		"NameI",
		"Position",
		"PersonId",
		"Points",
		"Rebounds",
		"Assists",
		"Fg3ptPct",
	}

	expectedVals := []string{
		"Rasheed Wallace",
		"F",
		"123456",
		"33",
		"18",
		"7",
		"39%",
	}

	cols, vals := getColsAndValues(player)

	if !reflect.DeepEqual(cols, expectedCols) {
		t.Errorf("expected: %v, got: %v", expectedCols, cols)
	}

	if !reflect.DeepEqual(vals, expectedVals) {
		t.Errorf("expected: %v, got: %v", expectedVals, vals)
	}
}

func TestExtractStatistics(t *testing.T) {

	stats := mockStats{
		Points:   33,
		Rebounds: 18,
		Assists:  7,
		Fg3ptPct: 0.388,
	}

	expectedKeys := []string{"Points", "Rebounds", "Assists", "Fg3ptPct"}
	expectedVals := []string{"33", "18", "7", "39%"}

	keys, vals := extractStatistics(stats)
	if !reflect.DeepEqual(expectedKeys, keys) {
		t.Errorf("expected keys: %v, got: %v", expectedKeys, keys)
	}

	if !reflect.DeepEqual(expectedVals, vals) {
		t.Errorf("expected values: %v, got: %v", expectedVals, vals)
	}
}

func TestExtractStatisticsZeroValues(t *testing.T) {
	stats := mockStats{
		Points:   0,
		Rebounds: 0,
		Assists:  0,
		Fg3ptPct: 0.000,
	}

	expectedVals := []string{"0", "0", "0", "0%"}

	_, vals := extractStatistics(stats)
	if !reflect.DeepEqual(expectedVals, vals) {
		t.Errorf("expected values: %v, got: %v", expectedVals, vals)
	}
}

func TestNewBoxScore(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	_, _, err := NewBoxScore("shittywok", WindowSize)
	if err == nil {
		t.Errorf("NewBoxScore() should have returned an error")
	}

}
