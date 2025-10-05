package benchmarking

func CocktailSort(nums []int) {
	start, end := 0, len(nums)
	swapped := true
	for swapped {
		swapped = false
		for i := start; i < end-1; i++ {
			if nums[i] > nums[i+1] {
				nums[i], nums[i+1] = nums[i+1], nums[i]
				swapped = true
			}
		}
		if !swapped {
			break
		}

		swapped = false
		end--

		for i := end - 1; i >= start; i-- {
			if nums[i] > nums[i+1] {
				nums[i], nums[i+1] = nums[i+1], nums[i]
				swapped = true
			}
		}
		start++
	}
}
