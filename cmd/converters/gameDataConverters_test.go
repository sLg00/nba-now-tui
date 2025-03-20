package converters

import (
	"github.com/sLg00/nba-now-tui/cmd/helpers"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"testing"
)

func mockUnmarshallDailyGameResults(_ string) (types.ResponseSet, error) {
	return types.ResponseSet{
		ResultSets: []types.ResultSet{
			{
				Name: "GameHeader", // Simulate the first set as metadata
				Headers: []string{
					"GAME_DATE_EST",
					"GAME_SEQUENCE",
					"GAME_ID",
					"GAME_STATUS_ID",
					"GAME_STATUS_TEXT",
					"GAMECODE",
					"HOME_TEAM_ID",
					"VISITOR_TEAM_ID",
				},
				RowSet: [][]interface{}{
					{"2024-12-01T00:00:00", float64(1), "001", float64(3), "Final", "20241201/NYKMIA", float64(1011), float64(1012)},
				},
			},
			{
				Name: "LineScore", // Second set as linescores
				Headers: []string{
					"GAME_ID", "TEAM_ID", "TEAM_ABBREVIATION", "PTS",
				},
				RowSet: [][]interface{}{
					{"001", float64(1011), "NYK", float64(101)},
					{"001", float64(1012), "MIA", float64(99)},
				},
			},
		},
	}, nil
}

func mockUnmarshallBoxScore(_ string) (types.ResponseSet, error) {
	return types.ResponseSet{
		BoxScore: types.BoxScore{
			GameID: "001",
			HomeTeam: types.BoxScoreTeam{
				TeamID:      101,
				TeamCity:    "New York",
				TeamName:    "Knicks",
				TeamTriCode: "NYK",
			},
			AwayTeam: types.BoxScoreTeam{
				TeamID:      102,
				TeamCity:    "Miami",
				TeamName:    "Heat",
				TeamTriCode: "MIA",
			},
		},
	}, nil
}

func TestPopulateDailyGameResults(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	bs, err := mockUnmarshallBoxScore("")
	if err != nil {
		t.Errorf("unmarshalling failed")
	}

	results, headers, err := PopulateDailyGameResults(bs)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %v", len(results))
	}

	if results[0].HomeTeamAbbreviation != "NYK" || results[0].HomeTeamPts != 101 {
		t.Errorf("Unexpected home team data %v", results[0])
	}

	if results[0].AwayTeamAbbreviation != "MIA" || results[0].AwayTeamPts != 99 {
		t.Errorf("Unexpected away team data %v", results[0])
	}

	if len(headers) != 4 {
		t.Errorf("Expected 4 headers, got %v", len(headers))
	}

}

func TestPopulateBoxScore_Success(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	bs, err := mockUnmarshallBoxScore("")
	if err != nil {
		t.Errorf("unmarshalling failed")
	}

	boxScore, err := PopulateBoxScore(bs)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if boxScore.GameID != "001" {
		t.Errorf("Expected GameID '001', got: %s", boxScore.GameID)
	}

	if boxScore.HomeTeam.TeamName != "Knicks" || boxScore.HomeTeam.TeamCity != "New York" {
		t.Errorf("Unexpected home team data: %+v", boxScore.HomeTeam)
	}

	if boxScore.AwayTeam.TeamName != "Heat" || boxScore.AwayTeam.TeamCity != "Miami" {
		t.Errorf("Unexpected away team data: %+v", boxScore.AwayTeam)
	}
}
