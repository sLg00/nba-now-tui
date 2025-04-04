package converters

import (
	"encoding/json"
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

// PopulateDailyGameResults extracts 'linescores' from the NBA API response for DailyScoreboard.
// Then subsequently converts the 'linescore' to GameResult objects, combining home and away team basic stats.
// The function could be cleanly split into two, but yolo for now.
func PopulateDailyGameResults(rs types.ResponseSet) (types.DailyGameResults, []string, error) {

	// ResultSets[1] is the "linescore" part of the response, which is used to create a gameCard
	if len(rs.ResultSets) < 2 || len(rs.ResultSets[1].Headers) == 0 {
		return nil, nil, fmt.Errorf("no headers or insufficient result sets in response")
	}

	gameHeaders := make(map[string]types.GameHeader)
	for _, row := range rs.ResultSets[0].RowSet {
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

	headers := rs.ResultSets[1].Headers
	rowSet := rs.ResultSets[1].RowSet

	var lineScores []types.LineScore
	for _, row := range rowSet {
		if len(row) != len(headers) {
			return nil, nil, fmt.Errorf("header row length does not match row length: %v != %v", len(headers), len(row))
		}

		gameData := make(map[string]interface{})
		for i, value := range row {
			gameData[headers[i]] = value
		}

		jsonData, err := json.Marshal(gameData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal json data: %v", err)
		}

		var lineScore types.LineScore
		err = json.Unmarshal(jsonData, &lineScore)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal json data: %v", err)
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
	cl, err := nbaAPI.NewClient().Loader.LoadDailyScoreboard()
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
	boxScore := rs.BoxScore

	return boxScore, nil
}
