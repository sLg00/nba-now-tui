package datamodels

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

// Stringer interface provides a method to convert the attributes of any type to strings.
// This is required to use types in the TUI, which requires objects to be strings
type Stringer interface {
	ToStringSlice() []string
}

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

// structToStringSlice is the core function that converts type attributes from Float64 and Int to String, using reflection
func structToStringSlice(obj any) []string {
	v := reflect.ValueOf(obj)
	var result []string

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		switch value.Kind() {
		case reflect.Float64:
			result = append(result, strconv.FormatFloat(value.Float(), 'f', 2, 64))
		case reflect.Int:
			result = append(result, strconv.Itoa(int(value.Int())))
		case reflect.String:
			result = append(result, value.String())
		}
	}
	return result
}

// ConvertToString takes the generic Stringer interface and does type conversion to string. All types that satisfy the
// Stringer interface can have their attributes converted to strings. For instance Player, Players, Team and Teams
func ConvertToString[T Stringer](objs []T) [][]string {
	var stringValues [][]string

	for _, obj := range objs {
		stringValues = append(stringValues, obj.ToStringSlice())
	}

	return stringValues
}
