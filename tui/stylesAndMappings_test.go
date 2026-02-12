package tui

import "testing"

func TestCalculatePageSize_SingleTable(t *testing.T) {
	// Terminal height 40, single table
	// baseOverhead = 4(margins) + 1(help) = 5
	// pageSize = 40 - 5 - 3(headerAndBorder) = 32
	result := calculatePageSize(40, 1)
	if result != 32 {
		t.Errorf("expected 32, got %d", result)
	}
}

func TestCalculatePageSize_DualTable(t *testing.T) {
	// Terminal height 40, dual tables
	// baseOverhead = 5, available = 35
	// pageSize = (35 - 6) / 2 = 14
	result := calculatePageSize(40, 2)
	if result != 14 {
		t.Errorf("expected 14, got %d", result)
	}
}

func TestCalculatePageSize_MinimumEnforced(t *testing.T) {
	// Tiny terminal should return minimum of 3
	result := calculatePageSize(10, 2)
	if result != 3 {
		t.Errorf("expected minimum 3, got %d", result)
	}
}
