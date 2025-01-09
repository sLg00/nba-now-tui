package client

import (
	"github.com/sLg00/nba-now-tui/cmd/helpers"
	"strconv"
	"strings"
	"testing"
)

func TestIdentifySeason(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	var seasonString string
	date, _ := GetDateArg()
	dateSplit := strings.Split(date, "-")
	year := dateSplit[0]
	month := dateSplit[1]

	monthInt, _ := strconv.Atoi(month)
	yearInt, _ := strconv.Atoi(year)
	previousYear := yearInt - 1
	nextYear := yearInt + 1

	if monthInt < 10 {
		dateStringPartOne := strconv.Itoa(previousYear)
		dateStringPartTwo := strconv.Itoa(yearInt)[2:]
		seasonString = dateStringPartOne + "-" + dateStringPartTwo

	}

	if monthInt >= 10 {
		dateStringPartOne := strconv.Itoa(yearInt)
		dateStringPartTwo := strconv.Itoa(nextYear)[2:]
		seasonString = dateStringPartOne + "-" + dateStringPartTwo

	}

	expectedResult := identifySeason()
	if expectedResult != seasonString {
		t.Errorf("expected: %s, got: %s", expectedResult, seasonString)
	}

}
