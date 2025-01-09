package tui

import (
	"github.com/sLg00/nba-now-tui/cmd/helpers"
	"testing"
)

func TestNewBoxScore(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	_, _, err := NewBoxScore("shittywok", WindowSize)
	if err == nil {
		t.Errorf("NewBoxScore() should have returned an error")
	}

}
