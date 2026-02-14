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

func TestCommonPlayerInfoParams_Endpoint(t *testing.T) {
	p := CommonPlayerInfoParams{PlayerID: "1628389"}
	if got := p.Endpoint(); got != "commonplayerinfo" {
		t.Errorf("Endpoint() = %s, want commonplayerinfo", got)
	}
}

func TestCommonPlayerInfoParams_Validate(t *testing.T) {
	p := CommonPlayerInfoParams{}
	if err := p.Validate(); err == nil {
		t.Error("Validate() expected error for missing PlayerID")
	}

	p.PlayerID = "1628389"
	if err := p.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}
}

func TestCommonPlayerInfoParams_ToValues(t *testing.T) {
	p := CommonPlayerInfoParams{PlayerID: "1628389"}
	v := p.ToValues()
	if got := v.Get("PlayerID"); got != "1628389" {
		t.Errorf("ToValues().Get(PlayerID) = %s, want 1628389", got)
	}
}

func TestPlayerProfileV2Params_Endpoint(t *testing.T) {
	p := PlayerProfileV2Params{PlayerID: "1628389", PerMode: "PerGame"}
	if got := p.Endpoint(); got != "playercareerstats" {
		t.Errorf("Endpoint() = %s, want playercareerstats", got)
	}
}

func TestPlayerProfileV2Params_Validate(t *testing.T) {
	p := PlayerProfileV2Params{}
	if err := p.Validate(); err == nil {
		t.Error("Validate() expected error for missing PlayerID")
	}

	p.PlayerID = "1628389"
	p.PerMode = "PerGame"
	if err := p.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}
}

func TestPlayerProfileV2Params_ToValues(t *testing.T) {
	p := PlayerProfileV2Params{PlayerID: "1628389", PerMode: "PerGame"}
	v := p.ToValues()
	if got := v.Get("PlayerID"); got != "1628389" {
		t.Errorf("ToValues().Get(PlayerID) = %s, want 1628389", got)
	}
	if got := v.Get("PerMode"); got != "PerGame" {
		t.Errorf("ToValues().Get(PerMode) = %s, want PerGame", got)
	}
}

func TestPlayerGameLogParams_Endpoint(t *testing.T) {
	p := PlayerGameLogParams{PlayerID: "1628389", Season: "2024-25", SeasonType: "Regular Season"}
	if got := p.Endpoint(); got != "playergamelog" {
		t.Errorf("Endpoint() = %s, want playergamelog", got)
	}
}

func TestPlayerGameLogParams_Validate(t *testing.T) {
	p := PlayerGameLogParams{}
	if err := p.Validate(); err == nil {
		t.Error("Validate() expected error for missing PlayerID")
	}

	p.PlayerID = "1628389"
	if err := p.Validate(); err == nil {
		t.Error("Validate() expected error for missing Season")
	}

	p.Season = "2024-25"
	p.SeasonType = "Regular Season"
	if err := p.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}
}

func TestPlayerGameLogParams_ToValues(t *testing.T) {
	p := PlayerGameLogParams{PlayerID: "1628389", Season: "2024-25", SeasonType: "Regular Season"}
	v := p.ToValues()
	if got := v.Get("PlayerID"); got != "1628389" {
		t.Errorf("ToValues().Get(PlayerID) = %s, want 1628389", got)
	}
	if got := v.Get("Season"); got != "2024-25" {
		t.Errorf("ToValues().Get(Season) = %s, want 2024-25", got)
	}
	if got := v.Get("SeasonType"); got != "Regular Season" {
		t.Errorf("ToValues().Get(SeasonType) = %s, want Regular Season", got)
	}
}
