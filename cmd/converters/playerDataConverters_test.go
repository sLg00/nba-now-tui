package converters

import (
	"github.com/sLg00/nba-now-tui/cmd/helpers"
	"slices"
	"testing"
)

func TestPlayerStatistics_ToStringSlice(t *testing.T) {
	testStats := PlayerStatistics{
		Minutes:                 "69:12",
		FieldGoalsMade:          30,
		FieldGoalsAttempted:     60,
		FieldGoalsPercentage:    50.0,
		ThreePointersMade:       5,
		ThreePointersAttempted:  20,
		ThreePointersPercentage: 25.0,
	}

	slice := testStats.ToStringSlice()
	if len(slice) == 0 {
		t.Errorf("testStats.ToStringSlice() returned empty slice")
	}

	if !slices.Contains(slice, "69:12") {
		t.Errorf("testStats.ToStringSlice() did not contain '69:12'")
	}
	if !slices.Contains(slice, "20") {
		t.Errorf("testStats.ToStringSlice() did not contain '20'")
	}

}

func TestBoxScorePlayer_ToStringSlice(t *testing.T) {
	player := BoxScorePlayer{
		PersonId:  123123,
		FirstName: "Domantas",
	}

	slice := player.ToStringSlice()
	if len(slice) == 0 {
		t.Errorf("player.ToStringSlice() returned empty slice")
	}

	if !slices.Contains(slice, "Domantas") {
		t.Errorf("player.ToStringSlice() did not contain 'Domantas'")
	}

	if !slices.Contains(slice, "123123") {
		t.Errorf("player.ToStringSlice() did not contain %s\n", slice)
	}
}

func TestPlayer_ToStringSlice(t *testing.T) {
	player := Player{
		PlayerID:   123123,
		PlayerName: "Jalen Brunson",
		TeamAbbr:   "NYK",
		PTS:        28.50,
	}

	slice := player.ToStringSlice()
	if len(slice) == 0 {
		t.Errorf("player.ToStringSlice() returned empty slice")
	}

	if !slices.Contains(slice, "Jalen Brunson") {
		t.Errorf("player.ToStringSlice() did not contain 'Jalen Brunson'")
	}

	if !slices.Contains(slice, "28.50") {
		t.Errorf("player.ToStringSlice() did not contain '28.50'")
	}
}

func TestPopulatePlayerStats(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	mockUnmarshall := func(path string) (ResponseSet, error) {
		return ResponseSet{
			ResultSet: ResultSet{
				Headers: []string{"PLAYER_ID", "PLAYER", "TEAM", "GP", "PTS"},
				RowSet: [][]interface{}{
					{1, "Domantas Sabonis", "SAC", 70, 28.5},
					{2, "Jalen Brunson", "NYK", 69, 27.9},
				},
			},
		}, nil
	}

	players, headers, err := PopulatePlayerStats(mockUnmarshall)

	if err != nil {
		t.Errorf("Expected no errors, got error %v", err)
	}

	if len(players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(players))
	}

	expectedHeaders := []string{"PLAYER_ID", "PLAYER", "TEAM", "GP", "PTS"}
	if len(headers) != len(expectedHeaders) {
		t.Errorf("Expected %d headers, got %d", len(expectedHeaders), len(headers))
	}

	if players[0].PlayerName != "Domantas Sabonis" {
		t.Errorf("Expected playerName Domantas Sabonis, got %s", players[0].PlayerName)
	}

	if players[1].TeamAbbr != "NYK" {
		t.Errorf("Expected team tag NYK, got %s", players[1].TeamAbbr)
	}
}
