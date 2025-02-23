package converters

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
)

// PopulateDailyGameResults extracts 'linescores' from the NBA API response for DailyScoreboard.
// Then subsequently converts the 'linescore' to GameResult objects, combining home and away team basic stats.
// The function could be cleanly split into two, but yolo for now.
func PopulateDailyGameResults(rs types.ResponseSet) (types.DailyGameResults, []string, error) {
	response := rs

	// ResultSets[1] is the "linescore" part of the response, which is used to create a gameCard
	if len(response.ResultSets) < 2 || len(response.ResultSets[1].Headers) == 0 {
		return nil, nil, fmt.Errorf("no headers or insufficient result sets in response")
	}

	gameHeaders := make(map[string]types.GameHeader)
	for _, row := range response.ResultSets[0].RowSet {
		if len(row) < 8 {
			err := fmt.Errorf("not enough header rows in response, expected at least 8 bytes, got %d", len(row))
			return nil, nil, err
		}
		gameHeader := types.GameHeader{
			GameID:        row[2].(string),
			GameStatusID:  int(row[3].(float64)),
			HomeTeamID:    int(row[6].(float64)),
			VisitorTeamID: int(row[7].(float64)),
		}
		gameHeaders[gameHeader.GameID] = gameHeader
	}

	headers := response.ResultSets[1].Headers
	rowSet := response.ResultSets[1].RowSet

	var lineScores []types.LineScore
	for _, row := range rowSet {
		if len(row) != len(headers) {
			log.Printf("header row length does not match row length: %v != %v", len(headers), len(row))
			return nil, nil, fmt.Errorf("header row length does not match row length: %v != %v", len(headers), len(row))
		}

		lineScore := types.LineScore{}
		for i, value := range row {
			header := headers[i]
			switch header {
			case "GAME_DATE_EST":
				if v, ok := value.(string); ok {
					lineScore.GameDateEst = v
				}
			case "GAME_ID":
				if v, ok := value.(string); ok {
					lineScore.GameID = v
				}
			case "TEAM_ID":
				if v, ok := value.(float64); ok {
					lineScore.TeamID = int(v)
				}
			case "TEAM_ABBREVIATION":
				if v, ok := value.(string); ok {
					lineScore.TeamAbbreviation = v
				}
			case "TEAM_CITY_NAME":
				if v, ok := value.(string); ok {
					lineScore.TeamCityName = v
				}
			case "PTS":
				if v, ok := value.(float64); ok {
					lineScore.Pts = int(v)
				}
			case "FG_PCT":
				if v, ok := value.(float64); ok {
					lineScore.FgPct = v
				}
			case "FT_PCT":
				if v, ok := value.(float64); ok {
					lineScore.FtPct = v
				}
			case "FG3_PCT":
				if v, ok := value.(float64); ok {
					lineScore.Fg3Pct = v
				}
			case "AST":
				if v, ok := value.(float64); ok {
					lineScore.Ast = int(v)
				}
			case "REB":
				if v, ok := value.(float64); ok {
					lineScore.Reb = int(v)
				}
			case "TOV":
				if v, ok := value.(float64); ok {
					lineScore.Tov = int(v)
				}
			}
		}
		lineScores = append(lineScores, lineScore)
	}

	var gameResults types.DailyGameResults
	for _, gh := range gameHeaders {
		result := types.GameResult{GameID: gh.GameID, GameStatusID: gh.GameStatusID}
		for _, ls := range lineScores {
			if ls.GameID == gh.GameID {
				if ls.TeamID == gh.HomeTeamID {
					result.HomeTeamID = ls.TeamID
					result.HomeTeamAbbreviation = ls.TeamAbbreviation
					result.HomeTeamPts = ls.Pts
				} else if ls.TeamID == gh.VisitorTeamID {
					result.AwayTeamID = ls.TeamID
					result.AwayTeamAbbreviation = ls.TeamAbbreviation
					result.AwayTeamPts = ls.Pts
				}
			}
		}
		gameResults = append(gameResults, result)
	}

	return gameResults, headers, nil
}

// CheckGameStatus takes a gameId and returns the game's status ID (1 - scheduled, 2 - live, 3 - final)
func CheckGameStatus(gameID string) (int, error) {
	cl, err := nbaAPI.NewNewClient().Loader.LoadDailyScoreboard()
	data, _, err := PopulateDailyGameResults(cl)
	if err != nil {
		return 0, fmt.Errorf("could not unmarshall json data: %v", err)
	}

	for _, game := range data {
		if game.GameID == gameID {
			return game.GameStatusID, nil
		}
	}
	return 0, fmt.Errorf("could not find game with id: %v", gameID)
}

// PopulateBoxScore takes a gameID (string) as input and returns the required structures to represent
// the game's box score in a TUI
func PopulateBoxScore(rs types.ResponseSet) (types.BoxScore, error) {
	response := rs
	boxScore := response.BoxScore

	return boxScore, nil
}
