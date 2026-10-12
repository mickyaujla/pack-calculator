package packing

import (
	"slices"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		items    int
		sizes    []int
		expected []Pack
	}{
		{"1 item", 1, DefaultSizes(), []Pack{{250, 1}}},
		{"exactly one pack", 250, DefaultSizes(), []Pack{{250, 1}}},
		{"one over a pack", 251, DefaultSizes(), []Pack{{500, 1}}},
		{"501 items", 501, DefaultSizes(), []Pack{{500, 1}, {250, 1}}},
		{"12001 items", 12001, DefaultSizes(), []Pack{{5000, 2}, {2000, 1}, {250, 1}}},
		{"greedy would send 46", 43, []int{6, 9, 20}, []Pack{{20, 1}, {9, 2}, {6, 1}}},
		{"large edge case", 500000, []int{23, 31, 53}, []Pack{{53, 9429}, {31, 7}, {23, 2}}},
		{"single pack size", 20, []int{7}, []Pack{{7, 3}}},
		{"unsorted sizes", 501, []int{1000, 250, 5000, 500, 2000}, []Pack{{500, 1}, {250, 1}}},
		{"duplicate sizes", 501, []int{250, 250, 500}, []Pack{{500, 1}, {250, 1}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Calculate(test.items, test.sizes)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !slices.Equal(got, test.expected) {
				t.Errorf("got %v, expected %v", got, test.expected)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	tests := []struct {
		name  string
		items int
		sizes []int
	}{
		{"zero items", 0, DefaultSizes()},
		{"negative items", -5, DefaultSizes()},
		{"too many items", MaxItems + 1, DefaultSizes()},
		{"no pack sizes", 10, []int{}},
		{"zero pack sizes", 10, []int{0, 250}},
		{"negative pack sizes", 10, []int{-250}},
		{"pack size too large", 10, []int{MaxItems + 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Calculate(test.items, test.sizes)

			if err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}
