package converters

import (
	"encoding/json"
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"reflect"
	"strings"
)

// structJSONHeaders extracts json tag names from a struct type, preserving field order.
// This ensures headers align with the output of ToStringSlice/ConvertToStringMatrix.
func structJSONHeaders(v interface{}) []string {
	t := reflect.TypeOf(v)
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" {
			continue
		}
		parts := strings.Split(tag, ",")
		headers = append(headers, parts[0])
	}
	return headers
}

// PopulatePlayerBio extracts the player bio from the commonplayerinfo API response.
// The response uses resultSets (plural) with the bio data in the first result set's single row.
func PopulatePlayerBio(rs types.ResponseSet) (types.PlayerBio, error) {
	if len(rs.ResultSets) == 0 || len(rs.ResultSets[0].RowSet) == 0 {
		return types.PlayerBio{}, fmt.Errorf("no player bio data found")
	}

	headers := rs.ResultSets[0].Headers
	row := rs.ResultSets[0].RowSet[0]

	if len(row) != len(headers) {
		return types.PlayerBio{}, fmt.Errorf("header length mismatch: %d vs %d", len(headers), len(row))
	}

	playerData := make(map[string]interface{})
	for i, value := range row {
		playerData[headers[i]] = value
	}

	jsonData, err := json.Marshal(playerData)
	if err != nil {
		return types.PlayerBio{}, fmt.Errorf("failed to marshal player bio: %v", err)
	}

	var bio types.PlayerBio
	if err = json.Unmarshal(jsonData, &bio); err != nil {
		return types.PlayerBio{}, fmt.Errorf("failed to unmarshal player bio: %v", err)
	}

	return bio, nil
}

// PopulateSeasonStats extracts career season stats from the playercareerstats API response.
// It looks for the result set named "SeasonTotalsRegularSeason".
func PopulateSeasonStats(rs types.ResponseSet) ([]types.SeasonStats, []string, error) {
	var targetSet *types.ResultSet
	for i := range rs.ResultSets {
		if rs.ResultSets[i].Name == "SeasonTotalsRegularSeason" {
			targetSet = &rs.ResultSets[i]
			break
		}
	}

	if targetSet == nil {
		return nil, nil, fmt.Errorf("SeasonTotalsRegularSeason result set not found")
	}

	apiHeaders := targetSet.Headers
	var stats []types.SeasonStats

	for _, row := range targetSet.RowSet {
		if len(row) != len(apiHeaders) {
			return nil, nil, fmt.Errorf("row length mismatch: %d vs %d", len(row), len(apiHeaders))
		}

		data := make(map[string]interface{})
		for i, value := range row {
			data[apiHeaders[i]] = value
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal season stats: %v", err)
		}

		var season types.SeasonStats
		if err = json.Unmarshal(jsonData, &season); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal season stats: %v", err)
		}

		stats = append(stats, season)
	}

	for i, j := 0, len(stats)-1; i < j; i, j = i+1, j-1 {
		stats[i], stats[j] = stats[j], stats[i]
	}

	structHeaders := structJSONHeaders(types.SeasonStats{})
	return stats, structHeaders, nil
}

// PopulateGameLog extracts the most recent games from the playergamelog API response.
// Returns at most 5 entries (the API returns games in reverse chronological order).
func PopulateGameLog(rs types.ResponseSet) ([]types.GameLogEntry, []string, error) {
	if len(rs.ResultSets) == 0 {
		return nil, nil, fmt.Errorf("no game log data found")
	}

	apiHeaders := rs.ResultSets[0].Headers
	rows := rs.ResultSets[0].RowSet

	limit := len(rows)
	if limit > 5 {
		limit = 5
	}

	var entries []types.GameLogEntry
	for _, row := range rows[:limit] {
		if len(row) != len(apiHeaders) {
			return nil, nil, fmt.Errorf("row length mismatch: %d vs %d", len(row), len(apiHeaders))
		}

		data := make(map[string]interface{})
		for i, value := range row {
			data[apiHeaders[i]] = value
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal game log: %v", err)
		}

		var entry types.GameLogEntry
		if err = json.Unmarshal(jsonData, &entry); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal game log: %v", err)
		}

		entries = append(entries, entry)
	}

	structHeaders := structJSONHeaders(types.GameLogEntry{})
	return entries, structHeaders, nil
}
