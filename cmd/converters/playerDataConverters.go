package converters

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
)

// PopulatePlayerStats maps the data to the Player struct
func PopulatePlayerStats(rs types.ResponseSet) (types.Players, []string, error) {
	response := rs

	// mapping as Player types
	headers := response.ResultSet.Headers
	var playerStats types.Players
	for _, row := range response.ResultSet.RowSet {
		if len(row) != len(headers) {
			err := fmt.Errorf("row length doesn't match headers length. %v != %v", len(row), len(headers))
			log.Println(err)
			return nil, nil, err
		}
		var player types.Player
		for i, value := range row {
			switch v := value.(type) {
			case float64:
				switch headers[i] {
				case "PLAYER_ID":
					player.PlayerID = int(v)
				case "RANK":
					player.Rank = int(v)
				case "GP":
					player.GamesPlayed = int(v)
				case "TEAM_ID":
					player.TeamID = int(v)
				case "MIN":
					player.Minutes = v
				case "FGM":
					player.FGM = v
				case "FGA":
					player.FGA = v
				case "FG_PCT":
					player.FGPCT = v
				case "FG3M":
					player.FG3PTM = v
				case "FG3A":
					player.FG3PTA = v
				case "FG3_PCT":
					player.FG3PTPCT = v
				case "FTM":
					player.FTM = v
				case "FTA":
					player.FTA = v
				case "FT_PCT":
					player.FTPCT = v
				case "OREB":
					player.OREB = v
				case "DREB":
					player.DREB = v
				case "REB":
					player.REB = v
				case "AST":
					player.AST = v
				case "STL":
					player.STL = v
				case "BLK":
					player.BLK = v
				case "TOV":
					player.TOV = v
				case "PTS":
					player.PTS = v
				case "EFF":
					player.EFF = v
				}
			case string:
				switch headers[i] {
				case "PLAYER":
					player.PlayerName = v
				case "TEAM":
					player.TeamAbbr = v
				}
			}
		}
		playerStats = append(playerStats, player)
	}
	return playerStats, headers, nil
}
