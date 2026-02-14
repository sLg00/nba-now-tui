package nbaAPI

import (
	"testing"
	"time"
)

func TestNewDateProvider_ReturnsEasternDate(t *testing.T) {
	dp := NewDateProvider()

	got, err := dp.GetCurrentDate()
	if err != nil {
		t.Fatalf("GetCurrentDate() unexpected error: %v", err)
	}

	eastern, _ := time.LoadLocation("America/New_York")
	want := time.Now().In(eastern).Format("2006-01-02")

	if got != want {
		t.Errorf("GetCurrentDate() = %s, want %s", got, want)
	}
}

func TestGetCurrentSeason_BeforeOctober(t *testing.T) {
	dp := &nbaDateProvider{date: "2025-03-15"}
	got := dp.GetCurrentSeason()
	if got != "2024-25" {
		t.Errorf("GetCurrentSeason() = %s, want 2024-25", got)
	}
}

func TestGetCurrentSeason_AfterOctober(t *testing.T) {
	dp := &nbaDateProvider{date: "2025-11-01"}
	got := dp.GetCurrentSeason()
	if got != "2025-26" {
		t.Errorf("GetCurrentSeason() = %s, want 2025-26", got)
	}
}
