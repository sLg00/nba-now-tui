package converters

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"strings"
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

// PopulateBoxScore returns the box score from a ResponseSet. For live games it prefers the
// CDN live data (rs.LiveGame), which carries real-time player stats; finished games fall back
// to the stats.nba.com response (rs.BoxScore).
func PopulateBoxScore(rs types.ResponseSet) (types.BoxScore, error) {
	if rs.LiveGame != nil {
		boxScore := types.BoxScore{
			GameID:   rs.LiveGame.GameID,
			HomeTeam: rs.LiveGame.HomeTeam,
			AwayTeam: rs.LiveGame.AwayTeam,
		}
		for i := range boxScore.HomeTeam.BoxScorePlayers {
			boxScore.HomeTeam.BoxScorePlayers[i].Statistics.Minutes = parseISO8601Duration(boxScore.HomeTeam.BoxScorePlayers[i].Statistics.Minutes)
		}
		for i := range boxScore.AwayTeam.BoxScorePlayers {
			boxScore.AwayTeam.BoxScorePlayers[i].Statistics.Minutes = parseISO8601Duration(boxScore.AwayTeam.BoxScorePlayers[i].Statistics.Minutes)
		}
		return boxScore, nil
	}
	return rs.BoxScore, nil
}

// parseISO8601Duration converts ISO 8601 duration strings (e.g. "PT13M42.30S") to "M:SS".
// Returns the input unchanged if it is not in that format.
func parseISO8601Duration(d string) string {
	d = strings.TrimPrefix(d, "PT")
	mIdx := strings.IndexByte(d, 'M')
	if mIdx < 0 {
		return d
	}
	minutes := d[:mIdx]
	rest := d[mIdx+1:]
	sIdx := strings.IndexByte(rest, '.')
	var seconds string
	if sIdx >= 0 {
		seconds = rest[:sIdx]
	} else {
		seconds = strings.TrimSuffix(rest, "S")
	}
	if len(seconds) == 1 {
		seconds = "0" + seconds
	}
	return minutes + ":" + seconds
}
