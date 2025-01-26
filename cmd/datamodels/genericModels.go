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
	GameDate     string `json:"GameDate"`
}

// ResultSet  is the object that represents the actual returned data structure or headers and rows
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
	BoxScore   BoxScore    `json:"boxScoreTraditional"`
}

// UnmarshallResponseJSON unmarshalls JSON from the appropriate JSON file.
// Takes string (full path to file) as an input and returns a ResponseSet struct
func UnmarshallResponseJSON(s string) (ResponseSet, error) {
	var response ResponseSet

	data, err := os.ReadFile(s)
	if err != nil {
		err = fmt.Errorf("error reading file %s", s)
		log.Println(err)
		return response, err
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		err = fmt.Errorf("error unmarshalling response from file: %s", s)
		log.Println(err)
		return ResponseSet{}, err
	}
	return response, nil
}

// structToStringSlice is the core function that converts type attributes from Float64 and Int to String,
// using reflection. It also looks for the "percentage" tag and if found, formats the values in a way that
// they will be presented as "xx%" in the UI
func structToStringSlice(obj any) []string {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	var result []string

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		field := t.Field(i)
		percentageTag := field.Tag.Get("percentage")

		switch value.Kind() {
		case reflect.Float64:
			if percentageTag == "true" {
				result = append(result, FloatToPercent(value.Float()))
			} else {
				result = append(result, strconv.FormatFloat(value.Float(), 'f', 2, 64))
			}
		case reflect.Int:
			result = append(result, strconv.Itoa(int(value.Int())))
		case reflect.String:
			result = append(result, value.String())
		default:
			log.Println("unhandled case")
			break
		}
	}
	return result
}

// ConvertToStringMatrix takes the generic Stringer interface and does type conversion to string. All types that satisfy the
// Stringer interface can have their attributes converted to strings. This function is specifically designed to
// take in a slice of structures and return a slice of a slice of strings
func ConvertToStringMatrix[T Stringer](objs []T) [][]string {
	var stringValues [][]string

	for _, obj := range objs {
		stringValues = append(stringValues, obj.ToStringSlice())
	}

	return stringValues
}

// ConvertToStringFlat takes the generic Stringer interface and does type conversion to string. All types that satisfy the
// Stringer interface can have their attributes converted to strings. This function is specifically designed to
// take in an object and return a slice of strings
func ConvertToStringFlat[T Stringer](obj T) []string {
	return obj.ToStringSlice()
}

// FloatToPercent converts a float to a string and formats it as a percentage representation
func FloatToPercent(f float64) string {
	return fmt.Sprintf("%.0f%%", f*100)
}
