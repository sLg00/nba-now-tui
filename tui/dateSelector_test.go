package tui

import (
	"strings"
	"testing"
	"time"
)

func TestNewDateSelector(t *testing.T) {
	ds := NewDateSelector("2025-02-14")
	if ds.date != "2025-02-14" {
		t.Errorf("expected date 2025-02-14, got %s", ds.date)
	}
	if ds.editing {
		t.Error("should start in navigation mode")
	}
	if ds.focused != true {
		t.Error("should start focused")
	}
}

func TestDateSelector_PreviousDay(t *testing.T) {
	ds := NewDateSelector("2025-02-14")
	ds.previousDay()
	if ds.date != "2025-02-13" {
		t.Errorf("expected 2025-02-13, got %s", ds.date)
	}
}

func TestDateSelector_NextDay_NotFuture(t *testing.T) {
	eastern, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(eastern).Format("2006-01-02")

	ds := NewDateSelector(today)
	ds.nextDay()
	if ds.date != today {
		t.Errorf("should not advance past today, got %s", ds.date)
	}
}

func TestDateSelector_NextDay_Normal(t *testing.T) {
	ds := NewDateSelector("2025-01-01")
	ds.nextDay()
	if ds.date != "2025-01-02" {
		t.Errorf("expected 2025-01-02, got %s", ds.date)
	}
}

func TestDateSelector_ValidateDate(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"2025-02-14", false},
		{"2025-13-01", true},
		{"not-a-date", true},
		{"2099-01-01", true},
	}

	for _, tt := range tests {
		err := validateDate(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateDate(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestDateSelector_View_NavigationMode(t *testing.T) {
	ds := NewDateSelector("2025-02-14")
	view := ds.View()
	if view == "" {
		t.Error("view should not be empty")
	}
	if !strings.Contains(view, "2025-02-14") {
		t.Error("view should contain the date")
	}
}
