package types

import (
	"reflect"
	"testing"
)

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
