package main

import (
	"errors"
	"testing"
)

func TestOrderCategory(t *testing.T) {
	tests := []struct {
		name     string
		amount   int
		category string
		wantErr  error
	}{
		{"negative", -1, "", errNegativeAmount},
		{"zero", 0, "small", nil},
		{"small", 7500, "small", nil},
		{"small upper boundary", 9999, "small", nil},
		{"regular lower boundary", 10000, "regular", nil},
		{"regular", 31000, "regular", nil},
		{"regular upper boundary", 49999, "regular", nil},
		{"large lower boundary", 50000, "large", nil},
		{"large", 100000, "large", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, err := orderCategory(tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("orderCategory(%d) error = %v, want %v", tt.amount, err, tt.wantErr)
			}
			if category != tt.category {
				t.Errorf("orderCategory(%d) category = %q, want %q", tt.amount, category, tt.category)
			}
		})
	}
}
