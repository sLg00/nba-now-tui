package converters

import (
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"testing"
)

func TestPopulatePlayerBio(t *testing.T) {
	rs := types.ResponseSet{
		ResultSets: []types.ResultSet{
			{
				Name:    "CommonPlayerInfo",
				Headers: []string{"FIRST_NAME", "LAST_NAME", "DISPLAY_FIRST_LAST", "TEAM_NAME", "TEAM_ABBREVIATION", "TEAM_CITY", "JERSEY", "POSITION", "HEIGHT", "WEIGHT", "COUNTRY", "SCHOOL", "DRAFT_YEAR", "DRAFT_ROUND", "DRAFT_NUMBER", "SEASON_EXP", "BIRTHDATE", "TEAM_ID"},
				RowSet: [][]interface{}{
					{"Bam", "Adebayo", "Bam Adebayo", "Heat", "MIA", "Miami", "13", "Center-Forward", "6-9", "255", "USA", "Kentucky", "2017", "1", "14", 8.0, "1997-07-18T00:00:00", 1610612748.0},
				},
			},
		},
	}

	bio, err := PopulatePlayerBio(rs)
	if err != nil {
		t.Fatalf("PopulatePlayerBio() error: %v", err)
	}

	if bio.DisplayName != "Bam Adebayo" {
		t.Errorf("expected DisplayName 'Bam Adebayo', got '%s'", bio.DisplayName)
	}
	if bio.TeamAbbreviation != "MIA" {
		t.Errorf("expected TeamAbbreviation 'MIA', got '%s'", bio.TeamAbbreviation)
	}
	if bio.JerseyNumber != "13" {
		t.Errorf("expected JerseyNumber '13', got '%s'", bio.JerseyNumber)
	}
	if bio.Position != "Center-Forward" {
		t.Errorf("expected Position 'Center-Forward', got '%s'", bio.Position)
	}
}

func TestPopulatePlayerBio_EmptyResultSets(t *testing.T) {
	rs := types.ResponseSet{}
	_, err := PopulatePlayerBio(rs)
	if err == nil {
		t.Error("expected error for empty ResultSets")
	}
}

func TestPopulateSeasonStats(t *testing.T) {
	rs := types.ResponseSet{
		ResultSets: []types.ResultSet{
			{
				Name:    "SeasonTotalsRegularSeason",
				Headers: []string{"SEASON_ID", "TEAM_ABBREVIATION", "GP", "MIN", "PTS", "REB", "AST", "STL", "BLK", "FG_PCT", "FG3_PCT", "FT_PCT"},
				RowSet: [][]interface{}{
					{"2022-23", "MIA", 75.0, 33.4, 20.4, 9.2, 3.2, 1.2, 0.8, 0.540, 0.188, 0.810},
					{"2023-24", "MIA", 71.0, 34.0, 19.3, 10.4, 3.9, 1.1, 0.9, 0.523, 0.152, 0.790},
					{"2024-25", "MIA", 55.0, 34.5, 18.4, 9.9, 2.8, 1.1, 0.8, 0.523, 0.200, 0.805},
				},
			},
		},
	}

	stats, headers, err := PopulateSeasonStats(rs)
	if err != nil {
		t.Fatalf("PopulateSeasonStats() error: %v", err)
	}

	if len(stats) != 3 {
		t.Fatalf("expected 3 seasons, got %d", len(stats))
	}

	if len(headers) == 0 {
		t.Error("expected non-empty headers")
	}

	if stats[0].SeasonID != "2022-23" {
		t.Errorf("expected first season '2022-23', got '%s'", stats[0].SeasonID)
	}

	if stats[2].TeamAbbr != "MIA" {
		t.Errorf("expected TeamAbbr 'MIA', got '%s'", stats[2].TeamAbbr)
	}
}

func TestPopulateSeasonStats_MissingResultSet(t *testing.T) {
	rs := types.ResponseSet{
		ResultSets: []types.ResultSet{
			{Name: "SomethingElse", Headers: []string{}, RowSet: nil},
		},
	}

	_, _, err := PopulateSeasonStats(rs)
	if err == nil {
		t.Error("expected error for missing SeasonTotalsRegularSeason")
	}
}

func TestPopulateGameLog(t *testing.T) {
	rs := types.ResponseSet{
		ResultSets: []types.ResultSet{
			{
				Name:    "PlayerGameLog",
				Headers: []string{"GAME_DATE", "MATCHUP", "WL", "MIN", "PTS", "REB", "AST", "STL", "BLK", "FG_PCT", "PLUS_MINUS"},
				RowSet: [][]interface{}{
					{"FEB 10, 2025", "MIA vs. BOS", "W", 36.0, 24.0, 10.0, 3.0, 1.0, 2.0, 0.550, 12.0},
					{"FEB 08, 2025", "MIA @ NYK", "L", 34.0, 18.0, 8.0, 4.0, 0.0, 1.0, 0.450, -5.0},
					{"FEB 06, 2025", "MIA vs. LAL", "W", 38.0, 28.0, 12.0, 2.0, 2.0, 3.0, 0.600, 15.0},
					{"FEB 04, 2025", "MIA @ CHI", "W", 32.0, 22.0, 9.0, 5.0, 1.0, 0.0, 0.480, 8.0},
					{"FEB 02, 2025", "MIA vs. ATL", "W", 30.0, 16.0, 7.0, 3.0, 1.0, 1.0, 0.500, 6.0},
					{"JAN 31, 2025", "MIA @ PHI", "L", 35.0, 20.0, 11.0, 2.0, 0.0, 2.0, 0.470, -3.0},
					{"JAN 29, 2025", "MIA vs. DEN", "W", 37.0, 25.0, 10.0, 4.0, 2.0, 1.0, 0.520, 10.0},
				},
			},
		},
	}

	entries, headers, err := PopulateGameLog(rs)
	if err != nil {
		t.Fatalf("PopulateGameLog() error: %v", err)
	}

	if len(entries) != 5 {
		t.Fatalf("expected 5 entries (capped), got %d", len(entries))
	}

	if len(headers) == 0 {
		t.Error("expected non-empty headers")
	}

	if entries[0].GameDate != "FEB 10, 2025" {
		t.Errorf("expected first game date 'FEB 10, 2025', got '%s'", entries[0].GameDate)
	}

	if entries[0].Matchup != "MIA vs. BOS" {
		t.Errorf("expected matchup 'MIA vs. BOS', got '%s'", entries[0].Matchup)
	}

	if entries[4].GameDate != "FEB 02, 2025" {
		t.Errorf("expected 5th entry date 'FEB 02, 2025', got '%s'", entries[4].GameDate)
	}
}

func TestPopulateGameLog_EmptyResultSets(t *testing.T) {
	rs := types.ResponseSet{}
	_, _, err := PopulateGameLog(rs)
	if err == nil {
		t.Error("expected error for empty ResultSets")
	}
}

func TestPopulateGameLog_FewerThan5Games(t *testing.T) {
	rs := types.ResponseSet{
		ResultSets: []types.ResultSet{
			{
				Name:    "PlayerGameLog",
				Headers: []string{"GAME_DATE", "MATCHUP", "WL", "MIN", "PTS", "REB", "AST", "STL", "BLK", "FG_PCT", "PLUS_MINUS"},
				RowSet: [][]interface{}{
					{"FEB 10, 2025", "MIA vs. BOS", "W", 36.0, 24.0, 10.0, 3.0, 1.0, 2.0, 0.550, 12.0},
					{"FEB 08, 2025", "MIA @ NYK", "L", 34.0, 18.0, 8.0, 4.0, 0.0, 1.0, 0.450, -5.0},
				},
			},
		},
	}

	entries, _, err := PopulateGameLog(rs)
	if err != nil {
		t.Fatalf("PopulateGameLog() error: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}
