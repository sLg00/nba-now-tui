package tui

import "testing"

func TestCalculatePageSize_SingleTable(t *testing.T) {
	// Terminal height 40, single table
	// Available = 40 - 4(margins) - 2(help) - 1(header) - 2(border) - 1(buffer) = 30
	result := calculatePageSize(40, 1)
	if result != 30 {
		t.Errorf("expected 30, got %d", result)
	}
}

func TestCalculatePageSize_DualTable(t *testing.T) {
	// Terminal height 40, dual tables
	// Available = 40 - 4(margins) - 2(help) - 2(labels) - 2(buffer) = 30
	// Per table = (30 - 6(overhead per table)) / 2 = 12
	result := calculatePageSize(40, 2)
	if result != 12 {
		t.Errorf("expected 12, got %d", result)
	}
}

func TestCalculatePageSize_MinimumEnforced(t *testing.T) {
	// Tiny terminal should return minimum of 3
	result := calculatePageSize(10, 2)
	if result != 3 {
		t.Errorf("expected minimum 3, got %d", result)
	}
}
