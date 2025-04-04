package converters

import (
	"encoding/json"
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

// PopulateTeamInfo maps the team info request result to the appropriate
// internal structs for further processing and use
func PopulateTeamInfo(rs types.ResponseSet) (types.TeamCommonInfo, []string, error) {
	headers := rs.ResultSets[0].Headers
	if len(headers) == 0 {
		return types.TeamCommonInfo{}, nil, fmt.Errorf("could not unmarshall team info: no %v returned", headers)
	}

	var currentTeam types.TeamCommonInfo
	for _, row := range rs.ResultSets[0].RowSet {
		if len(row) != len(headers) {
			err := fmt.Errorf("header length mismatch with row: %v vs %v", len(headers), len(row))
			return types.TeamCommonInfo{}, nil, err
		}

		teamData := make(map[string]interface{})
		for i, value := range row {
			teamData[headers[i]] = value
		}

		jsonData, err := json.Marshal(teamData)
		if err != nil {
			return types.TeamCommonInfo{}, nil, fmt.Errorf("could not marshall team info: %v", err)
		}

		err = json.Unmarshal(jsonData, &currentTeam)
		if err != nil {
			return types.TeamCommonInfo{}, nil, fmt.Errorf("could not unmarshall team info: %v", err)
		}
	}

	return currentTeam, headers, nil
}

// PopulateTeamStats maps the data to Team struct to display season standings. Since headers and rows are in distinct
// json structures when calling NBA API, the function merges these into a map, marshalls and then unmarshalls them.
func PopulateTeamStats(rs types.ResponseSet) (types.Teams, []string, error) {
	headers := rs.ResultSets[0].Headers
	if len(headers) == 0 {
		return types.Teams{}, nil, fmt.Errorf("could not unmarshall team stats: no heads returned")
	}

	var teams types.Teams
	for _, row := range rs.ResultSets[0].RowSet {
		if len(row) != len(headers) {
			return nil, nil, fmt.Errorf("header length mismatch with row: %v vs %v", len(headers), len(row))
		}

		teamData := make(map[string]interface{})
		for i, value := range row {
			teamData[headers[i]] = value
		}

		jsonData, err := json.Marshal(teamData)
		if err != nil {
			return nil, nil, fmt.Errorf("could not marshall team stats: %v", err)
		}

		var team types.Team
		err = json.Unmarshal(jsonData, &team)
		if err != nil {
			return nil, nil, fmt.Errorf("could not unmarshall team stats: %v", err)
		}

		teams = append(teams, team)
	}

	return teams, headers, nil
}
