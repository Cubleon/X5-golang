package main

import (
	"errors"
	"testing"
)

func TestCalculateStats(t *testing.T) {
	tests := []struct {
		name    string
		amounts []int
		min     int
		max     int
		sum     int
		average float64
		wantErr error
	}{
		{"example", []int{12000, 7500, 31000, 4500}, 4500, 31000, 55000, 13750, nil},
		{"single order", []int{7500}, 7500, 7500, 7500, 7500, nil},
		{"equal amounts", []int{10, 10, 10}, 10, 10, 30, 10, nil},
		{"fractional average", []int{1, 2}, 1, 2, 3, 1.5, nil},
		{"ascending", []int{1, 2, 3}, 1, 3, 6, 2, nil},
		{"descending", []int{3, 2, 1}, 1, 3, 6, 2, nil},
		{"zero amounts", []int{0, 0}, 0, 0, 0, 0, nil},
		{"negative amounts", []int{-5, -1, -3}, -5, -1, -9, -3, nil},
		{"mixed amounts", []int{-5, 0, 8}, -5, 8, 3, 1, nil},
		{"empty slice", []int{}, 0, 0, 0, 0, errEmptyOrders},
		{"nil slice", nil, 0, 0, 0, 0, errEmptyOrders},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max, sum, average, err := calculateStats(tt.amounts)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("calculateStats(%v) error = %v, want %v", tt.amounts, err, tt.wantErr)
			}
			if min != tt.min || max != tt.max || sum != tt.sum || average != tt.average {
				t.Errorf("calculateStats(%v) = (%d, %d, %d, %g), want (%d, %d, %d, %g)",
					tt.amounts, min, max, sum, average, tt.min, tt.max, tt.sum, tt.average)
			}
		})
	}
}
