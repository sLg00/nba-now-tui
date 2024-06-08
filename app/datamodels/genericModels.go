package datamodels

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
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

// ConvertToString is a helper function that takes any struct type and converts its values to strings
// it is required since Bubble Tea uses strings to draw the UI
func ConvertToString[T any](obj []T) [][]string {
	var statsString [][]string

	for _, row := range obj {
		var instance []string

		v := reflect.ValueOf(row)

		for i := 0; i < v.NumField(); i++ {
			value := v.Field(i)
			switch value.Interface().(type) {
			case float64:
				instance = append(instance, strconv.FormatFloat(value.Float(), 'f', 2, 64))
			case int:
				instance = append(instance, strconv.Itoa(int(value.Int())))
			case string:
				instance = append(instance, value.String())
			}
		}
		statsString = append(statsString, instance)
	}
	return statsString
}
