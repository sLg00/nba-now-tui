package converters

import (
	"github.com/sLg00/nba-now-tui/cmd/helpers"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"testing"
)

func mockUnmarshallDailyGameResults(_ string) (types.ResponseSet, error) {
	return types.ResponseSet{
		Scoreboard: &types.ScoreboardV3Data{
			Games: []types.ScoreboardV3Game{
				{
					GameID:     "001",
					GameStatus: 3,
					HomeTeam: types.ScoreboardV3Team{
						TeamID:      1011,
						TeamTricode: "NYK",
						Score:       101,
					},
					AwayTeam: types.ScoreboardV3Team{
						TeamID:      1012,
						TeamTricode: "MIA",
						Score:       99,
					},
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

	bs, err := mockUnmarshallDailyGameResults("")
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

	if len(headers) != 6 {
		t.Errorf("Expected 6 headers, got %v", len(headers))
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
