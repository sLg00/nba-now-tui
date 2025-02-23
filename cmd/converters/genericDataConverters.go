package converters

import (
	"encoding/json"
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
	"os"
)

// UnmarshallResponseJSON unmarshalls JSON from the appropriate JSON file.
// Takes string (full path to file) as an input and returns a ResponseSet struct
func UnmarshallResponseJSON(s string) (types.ResponseSet, error) {
	var response types.ResponseSet

	data, err := os.ReadFile(s)
	if err != nil {
		err = fmt.Errorf("error reading file %s", s)
		log.Println(err)
		return response, err
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		err = fmt.Errorf("error unmarshalling response from file: %s, error %s ", s, err)
		log.Println(err)
		return types.ResponseSet{}, err
	}
	return response, nil
}
