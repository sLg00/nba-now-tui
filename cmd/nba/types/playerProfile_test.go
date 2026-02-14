package types

import (
	"reflect"
	"slices"
	"testing"
)

func TestPlayerBio_ToStringSlice(t *testing.T) {
	bio := PlayerBio{
		DisplayName:      "Bam Adebayo",
		TeamName:         "Heat",
		TeamAbbreviation: "MIA",
		JerseyNumber:     "13",
		Position:         "Center-Forward",
		Height:           "6-9",
		Weight:           "255",
		Country:          "USA",
		School:           "Kentucky",
		SeasonExp:        8,
		TeamID:           1610612748,
	}

	slice := bio.ToStringSlice()
	if len(slice) == 0 {
		t.Fatal("PlayerBio.ToStringSlice() returned empty slice")
	}
	if !slices.Contains(slice, "Bam Adebayo") {
		t.Errorf("ToStringSlice() did not contain 'Bam Adebayo', got %v", slice)
	}
	if !slices.Contains(slice, "MIA") {
		t.Errorf("ToStringSlice() did not contain 'MIA', got %v", slice)
	}
}

func TestSeasonStats_ToStringSlice(t *testing.T) {
	ss := SeasonStats{
		SeasonID: "2024-25",
		TeamAbbr: "MIA",
		GP:       55,
		PTS:      18.4,
		REB:      9.9,
		AST:      2.8,
		FGPCT:    0.523,
	}

	slice := ss.ToStringSlice()
	if len(slice) == 0 {
		t.Fatal("SeasonStats.ToStringSlice() returned empty slice")
	}
	if !slices.Contains(slice, "2024-25") {
		t.Errorf("ToStringSlice() did not contain '2024-25', got %v", slice)
	}
	if !slices.Contains(slice, "55") {
		t.Errorf("ToStringSlice() did not contain '55' for GP, got %v", slice)
	}
	if !slices.Contains(slice, "52%") {
		t.Errorf("ToStringSlice() did not contain '52%%' for FGPCT, got %v", slice)
	}
}

func TestGameLogEntry_ToStringSlice(t *testing.T) {
	entry := GameLogEntry{
		GameDate:  "FEB 10, 2025",
		Matchup:   "MIA vs. BOS",
		WL:        "W",
		MIN:       36,
		PTS:       24,
		REB:       10,
		AST:       3,
		STL:       1,
		BLK:       2,
		FGPCT:     0.55,
		PlusMinus: 12.0,
	}

	slice := entry.ToStringSlice()
	if len(slice) == 0 {
		t.Fatal("GameLogEntry.ToStringSlice() returned empty slice")
	}
	if !slices.Contains(slice, "MIA vs. BOS") {
		t.Errorf("ToStringSlice() did not contain 'MIA vs. BOS', got %v", slice)
	}
	if !slices.Contains(slice, "24") {
		t.Errorf("ToStringSlice() did not contain '24' for PTS, got %v", slice)
	}
}

func TestSeasonStats_StructTags(t *testing.T) {
	typ := reflect.TypeOf(SeasonStats{})

	tests := []struct {
		field   string
		tagKey  string
		tagVal  string
	}{
		{"FGPCT", "percentage", "true"},
		{"FG3PCT", "percentage", "true"},
		{"FTPCT", "percentage", "true"},
		{"SeasonID", "display", "Season"},
		{"PTS", "display", "PPG"},
		{"GP", "isVisible", "true"},
	}

	for _, tt := range tests {
		f, ok := typ.FieldByName(tt.field)
		if !ok {
			t.Errorf("field %s not found on SeasonStats", tt.field)
			continue
		}
		got := f.Tag.Get(tt.tagKey)
		if got != tt.tagVal {
			t.Errorf("SeasonStats.%s tag %s = %q, want %q", tt.field, tt.tagKey, got, tt.tagVal)
		}
	}
}
