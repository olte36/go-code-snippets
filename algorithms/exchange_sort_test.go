package algorithms

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchangeSort(t *testing.T) {
	cases := []struct {
		name  string
		input []int
	}{
		{"empty", []int{}},
		{"single element", []int{1}},
		{"already sorted", []int{1, 2, 3, 4, 5}},
		{"reverse order", []int{5, 4, 3, 2, 1}},
		{"duplicates", []int{3, 1, 2, 3, 2, 1}},
		{"negatives and positives", []int{-1, 5, 0, -3, 2}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			want := make([]int, len(tc.input))
			copy(want, tc.input)

			slices.Sort(want)

			got := make([]int, len(tc.input))
			copy(got, tc.input)

			ExchangeSort(got)

			assert.EqualValues(t, want, got)
		})
	}
}
