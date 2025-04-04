package converters

import (
	"encoding/json"
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
)

// PopulatePlayerStats maps the data to the Player struct
func PopulatePlayerStats(rs types.ResponseSet) (types.Players, []string, error) {
	// mapping as Player types
	headers := rs.ResultSet.Headers
	var playerStats types.Players
	for _, row := range rs.ResultSet.RowSet {
		if len(row) != len(headers) {
			log.Println("len(row) != len(headers):", len(row), len(headers))
			return nil, nil, fmt.Errorf("row length doesn't match headers length. %v != %v", len(row), len(headers))
		}

		playerData := make(map[string]interface{})
		for i, value := range row {
			playerData[headers[i]] = value
		}

		jsonData, err := json.Marshal(playerData)
		if err != nil {
			log.Println(err)
			return nil, nil, fmt.Errorf("failed to marshal player data. %v", err)
		}

		var player types.Player
		err = json.Unmarshal(jsonData, &player)
		if err != nil {
			log.Println(err)
			return nil, nil, fmt.Errorf("failed to unmarshal player data. %v", err)
		}

		playerStats = append(playerStats, player)
	}
	return playerStats, headers, nil
}
