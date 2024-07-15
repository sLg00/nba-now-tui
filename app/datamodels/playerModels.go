package datamodels

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/app/internal/client"
)

// Player struct represents a player row with their current statistical averages based on the input parameters
// Can be totals, per game averages, per 48 minutes etc.
type Player struct {
	PlayerID    int     `json:"PLAYER_ID"`
	Rank        int     `json:"RANK"`
	PlayerName  string  `json:"PLAYER"`
	TeamID      int     `json:"TEAM_ID"`
	TeamAbbr    string  `json:"TEAM"`
	GamesPlayed int     `json:"GP"`
	Minutes     float64 `json:"MIN"`
	FGM         float64 `json:"FGM"`
	FGA         float64 `json:"FGA"`
	FGPCT       float64 `json:"FG_PCT"`
	FG3PTM      float64 `json:"FG3PTM"`
	FG3PTA      float64 `json:"FG3PTA"`
	FG3PTPCT    float64 `json:"FG3PT_PCT"`
	FTM         float64 `json:"FTM"`
	FTA         float64 `json:"FTA"`
	FTPCT       float64 `json:"FT_PCT"`
	OREB        float64 `json:"OREB"`
	DREB        float64 `json:"DREB"`
	REB         float64 `json:"REB"`
	AST         float64 `json:"AST"`
	STL         float64 `json:"STL"`
	BLK         float64 `json:"BLK"`
	TOV         float64 `json:"TOV"`
	PTS         float64 `json:"PTS"`
	EFF         float64 `json:"EFF"`
}

// ToStringSlice is a method on the Player type that enables the attributes of the type to be converted to strings
func (p Player) ToStringSlice() []string {
	return structToStringSlice(p)
}

type Players []Player

// ToStringSlice is a method on the Players type that enables the attributes of the type to be converted to strings
func (ps Players) ToStringSlice() []string {
	return structToStringSlice(ps)
}

// PopulatePlayerStats maps the data to the Player struct
func PopulatePlayerStats() (Players, []string, error) {
	pc := client.NewClient().InstantiatePaths().LLFullPath()
	response, err := unmarshallResponseJSON(pc)
	if err != nil {
		fmt.Println("Could not unmarshall json data:", err)
	}
	// mapping as Player types. the code is ugly with the switch statements, but it works
	headers := response.ResultSet.Headers
	var playerStats Players
	for _, row := range response.ResultSet.RowSet {
		if len(row) != len(headers) {
			fmt.Println("Error: Row length doesn't match headers length")
			return nil, nil, err
		}
		var player Player
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
