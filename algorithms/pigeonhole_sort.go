package algorithms

import "slices"

// PigeonholeSort sorts the provided slice of ints in-place in ascending order
// using the pigeonhole sort algorithm.
//
// The implementation finds the minimum and maximum values to determine the
// pigeonhole range, allocates a counting slice of length max-min+1, counts
// occurrences of each key, and then rewrites the input slice from those
// counts. Negative values are handled by offsetting indices by the minimum
// value.
//
// Complexity: time O(n + r) and space O(r), where n is len(arr) and r is the
// value range (max - min). Because extra memory is proportional to the value
// range, this algorithm is appropriate only when the range is reasonably small
// relative to the number of elements.
func PigeonholeSort(arr []int) {
	min := slices.Min(arr)
	max := slices.Max(arr)
	// a slice with the length of the range of values in arr
	holes := make([]int, max-min+1)

	for i := range arr {
		holes[arr[i]-min]++
	}

	var arrInd int
	for i := range holes {
		for holes[i] > 0 {
			arr[arrInd] = i + min
			arrInd++
			holes[i]--
		}
	}
}
