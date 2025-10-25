package converters

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

// PopulateDailyGameResults extracts 'linescores' from the NBA API response for DailyScoreboard.
// Then subsequently converts the 'linescore' to GameResult objects, combining home and away team basic stats.
// The function could be cleanly split into two, but yolo for now.
func PopulateDailyGameResults(rs types.ResponseSet) (types.DailyGameResults, []string, error) {
	if rs.Scoreboard == nil || len(rs.Scoreboard.Games) == 0 {
		return nil, nil, fmt.Errorf("no games found in scoreboard response")
	}

	var gameResults types.DailyGameResults
	headers := []string{"GameID", "HomeTeam", "HomeScore", "AwayTeam", "AwayScore", "Status"}

	for _, game := range rs.Scoreboard.Games {
		result := types.GameResult{
			GameID:               game.GameID,
			GameStatusID:         game.GameStatus,
			HomeTeamID:           game.HomeTeam.TeamID,
			HomeTeamAbbreviation: game.HomeTeam.TeamTricode,
			HomeTeamName:         game.HomeTeam.TeamName,
			HomeTeamPts:          game.HomeTeam.Score,
			AwayTeamID:           game.AwayTeam.TeamID,
			AwayTeamAbbreviation: game.AwayTeam.TeamTricode,
			AwayTeamName:         game.AwayTeam.TeamName,
			AwayTeamPts:          game.AwayTeam.Score,
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
