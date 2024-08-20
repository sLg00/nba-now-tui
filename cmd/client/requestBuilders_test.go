package client

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestIdentifySeason(t *testing.T) {
	var s string
	cyms := strings.Split(time.Now().Format("2006-01"), "-")
	year := cyms[0]
	month := cyms[1]
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)
	lastYear := yearInt - 1
	nextYear := yearInt + 1

	if monthInt >= 10 {
		p1 := strconv.Itoa(yearInt)
		p2 := strconv.Itoa(nextYear)[2:]
		s = p1 + "-" + p2

	} else {
		p1 := strconv.Itoa(lastYear)
		p2 := strconv.Itoa(yearInt)[2:]
		s = p1 + "-" + p2

		expectedResult := identifySeason()
		if expectedResult != s {
			t.Errorf("expected: %s, got: %s", expectedResult, s)
		}
	}

}
