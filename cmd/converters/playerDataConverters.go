package converters

import (
	"encoding/json"
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

// PopulatePlayerStats maps the data to the Player struct
func PopulatePlayerStats(rs types.ResponseSet) (types.Players, []string, error) {
	// mapping as Player types
	headers := rs.ResultSet.Headers
	var playerStats types.Players
	for _, row := range rs.ResultSet.RowSet {
		if len(row) != len(headers) {
			return nil, nil, fmt.Errorf("row length doesn't match headers length. %v != %v", len(row), len(headers))
		}

		playerData := make(map[string]interface{})
		for i, value := range row {
			playerData[headers[i]] = value
		}

		jsonData, err := json.Marshal(playerData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal player data. %v", err)
		}

		var player types.Player
		err = json.Unmarshal(jsonData, &player)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal player data. %v", err)
		}

		playerStats = append(playerStats, player)
	}
	return playerStats, headers, nil
}

func PopulatePlayerIndex(rs types.ResponseSet) (types.IndexPlayers, []string, error) {
	headers := rs.ResultSets[0].Headers

	var playerIndex types.IndexPlayers
	for _, row := range rs.ResultSets[0].RowSet {
		if len(row) != len(headers) {
			return nil, nil, fmt.Errorf("row length doesn't match headers length. %v", len(row))
		}

		playerData := make(map[string]interface{})
		for i, value := range row {
			playerData[headers[i]] = value
		}

		jsonData, err := json.Marshal(playerData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal player data. %v", err)
		}

		var player types.IndexPlayer
		err = json.Unmarshal(jsonData, &player)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal player data. %v", err)
		}

		playerIndex = append(playerIndex, player)
	}

	return playerIndex, headers, nil
}
