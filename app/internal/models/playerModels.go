package models

import (
	"encoding/json"
	"fmt"
	"github.com/sLg00/nba-now-tui/app/internal/client"
)

// Parameters struct represents the parameters headers returned with the JSON response from the stats API
type Parameters struct {
	LeagueID     string `json:"LeagueID"`
	PerMode      string `json:"PerMode"`
	StatCategory string `json:"StatCategory"`
	Season       string `json:"Season"`
	SeasonType   string `json:"SeasonType"`
	Scope        string `json:"Scope"`
	ActiveFlag   string `json:"ActiveFlag"`
}

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

// ResultSet  is the object that represents the actual returned data structure or headers and rows.
type ResultSet struct {
	Name    string          `json:"name"`
	Headers []string        `json:"headers"`
	RowSet  [][]interface{} `json:"rowSet"`
}

// ResponseSet is the object to which the incoming JSON is unmarshalled
type ResponseSet struct {
	Resource   string     `json:"resource"`
	Parameters Parameters `json:"parameters"`
	ResultSet  ResultSet  `json:"resultSet"`
}

// PopulatePlayerStats unmarshalls the returned JSON and maps the data to the the Player struct
func PopulatePlayerStats() {

	//unmarshalling
	var response ResponseSet
	jsonData := client.InitiateClient()

	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	// mapping as Player types. the code is ugly with the switch statements, but it works
	var playerStats []Player
	for _, row := range response.ResultSet.RowSet {
		if len(row) != len(response.ResultSet.Headers) {
			fmt.Println("Error: Row length doesn't match headers length")
			return
		}
		var player Player
		for i, value := range row {
			switch v := value.(type) {
			case float64:
				switch response.ResultSet.Headers[i] {
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
					player.FTA = v
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
				switch response.ResultSet.Headers[i] {
				case "PLAYER":
					player.PlayerName = v
				case "TEAM":
					player.TeamAbbr = v
				}
			}
		}
		// playerStats contains all player data in a slice. TODO: write the results to a file with the current date!
		playerStats = append(playerStats, player)
	}
}
