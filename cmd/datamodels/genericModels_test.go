package datamodels

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUnmarshallRequestJSON(t *testing.T) {
	//temporarily using files already on the filesys. TODO: refactor so that files are set up and torn down specifically
	//for the tests
	mockFile := "/home/omar_t/.config/nba-tui/2024-07-10_ll"
	mockResponse, _ := UnmarshallResponseJSON(mockFile)
	mockFileWithNoData := "/home/omar_t/.config/nba-tui/2024-07-07_ll"
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
			result, _ := UnmarshallResponseJSON(test.input)
			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("Wanted: %v, Got: %v, Returned err: %v", test.want, result, test.err)
			}
		})
	}
}

func TestStructToStringSlice(t *testing.T) {

	type args struct {
		intArg    int
		floatArg  float64
		stringArg string
	}

	testArgs := args{
		intArg:    1,
		floatArg:  1.00,
		stringArg: "one",
	}

	var tests = []struct {
		name  string
		input interface{}
		want  []string
	}{
		{"coveredtypes", testArgs, []string{"1", "1.00", "one"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := structToStringSlice(test.input)
			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("Wanted: %v, Got: %v", test.want, result)
			}
		})
	}
}

func TestFloatToPercent(t *testing.T) {
	value := 0.1
	result := FloatToPercent(value)
	if result != "10%" {
		t.Errorf("Wanted: %v, Got: %v", "1%", result)
	}
}
