package datamodels

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Parameters struct represents the parameters headers returned with the JSON response from the stats API
type Parameters struct {
	LeagueID     string `json:"LeagueID"`
	PerMode      string `json:"PerMode"`
	StatCategory string `json:"StatCategory"`
	Season       string `json:"Season"`
	SeasonYear   string `json:"SeasonYear"`
	SeasonType   string `json:"SeasonType"`
	Scope        string `json:"Scope"`
	ActiveFlag   string `json:"ActiveFlag"`
}

// ResultSet  is the object that represents the actual returned data structure or headers and rows.SeasonYear
type ResultSet struct {
	Name    string          `json:"name"`
	Headers []string        `json:"headers"`
	RowSet  [][]interface{} `json:"rowSet"`
}

// ResponseSet is the object to which the incoming JSON is unmarshalled
type ResponseSet struct {
	Resource   string      `json:"resource"`
	Parameters Parameters  `json:"parameters"`
	ResultSet  ResultSet   `json:"resultSet"`
	ResultSets []ResultSet `json:"resultSets"`
}

// unmarshallResponseJSON unmarshalls JSON from the appropriate JSON file. Takes string (full path to file) as an input
func unmarshallResponseJSON(s string) (ResponseSet, error) {
	var response ResponseSet

	data, err := os.ReadFile(s)
	if err != nil {
		log.Println("file could not be found")
		return response, err
	}
	jsonData := data
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		fmt.Println("err:", err)
		return ResponseSet{}, err
	}
	return response, nil
}
