package algorithms

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPigeonholeSort(t *testing.T) {
	cases := []struct {
		name string
		in   []int
	}{
		{"single element", []int{42}},
		{"already sorted", []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}},
		{"duplicates", []int{3, 1, 2, 3, 2, 1, 3}},
		{"negatives", []int{-2, 0, -5, 3, -2}},
		{"mixed range", []int{1000, -1000, 0, 500, -500}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Clone(tc.in)
			PigeonholeSort(got)

			want := slices.Clone(tc.in)
			slices.Sort(want)

			assert.EqualValues(t, want, got)
		})
	}
}
