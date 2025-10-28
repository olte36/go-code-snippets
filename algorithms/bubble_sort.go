package algorithms

func BubbleSort(arr []int) {
	n := len(arr)
	for n > 1 {
		var swapped bool
		for i := 0; i < n-1; i++ {
			if arr[i] > arr[i+1] {
				arr[i], arr[i+1] = arr[i+1], arr[i]
				swapped = true
			}
		}
		// if no swaps occurred, the array is sorted
		if !swapped {
			break
		}
		// the largest element is now at the end, reduce the range
		n--
	}
}
