package benchmarking

import "testing"

func TestCocktailSort(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{
			name:     "unsorted",
			nums:     []int{4, 2, 5, 1, 6, 3},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "desc sorted",
			nums:     []int{6, 5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			CocktailSort(tc.nums)
			if !equalSlices(tc.nums, tc.expected) {
				t.Errorf("slice should be sorted to %v but got %v", tc.expected, tc.nums)
			}
		})
	}
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
