package tui

import "testing"

func TestCalculatePageSize_SingleTable(t *testing.T) {
	// Terminal height 40, single table
	// baseOverhead = 4(margins) + 1(help) = 5
	// pageSize = 40 - 5 - 4(headerAndBorder) = 31
	result := calculatePageSize(40, 1)
	if result != 31 {
		t.Errorf("expected 31, got %d", result)
	}
}

func TestCalculatePageSize_DualTable(t *testing.T) {
	// Terminal height 40, dual tables
	// baseOverhead = 5, available = 35
	// pageSize = (35 - 8) / 2 = 13
	result := calculatePageSize(40, 2)
	if result != 13 {
		t.Errorf("expected 13, got %d", result)
	}
}

func TestCalculatePageSize_MinimumEnforced(t *testing.T) {
	// Tiny terminal should return minimum of 3
	result := calculatePageSize(10, 2)
	if result != 3 {
		t.Errorf("expected minimum 3, got %d", result)
	}
}
