package datamodels

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUnmarshallRequestJSON(t *testing.T) {

	mockFile := "/home/omar_t/.config/nba-tui/2024-07-10_ll" //somefile that has the expected structure
	mockResponse, _ := unmarshallResponseJSON(mockFile)
	mockFileWithNoData := "/home/omar_t/.config/nba-tui/2024-07-07_ll" //somefile that has an unexpected structure
	missingFile := "/tmp/filenotfound.conf"

	var tests = []struct {
		name  string
		input string
		want  ResponseSet
		err   error
	}{
		{"correct result", mockFile, mockResponse, nil},
		{"faulty data", mockFileWithNoData, ResponseSet{}, nil},
		{"missing file", missingFile, ResponseSet{},
			fmt.Errorf("error reading file %s", missingFile)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, _ := unmarshallResponseJSON(test.input)
			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("Wanted: %v, Got: %v, Returned err: %v", test.want, result, test.err)
			}
		})
	}

}
