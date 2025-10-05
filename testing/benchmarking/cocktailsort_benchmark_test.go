package benchmarking

import (
	"math/rand"
	"testing"
)

func BenchmarkCocktailSortRandom(b *testing.B) {
	testCases := []struct {
		name string
		nums []int
	}{
		{
			name: "10 elements",
			nums: randomSlice(10),
		},
		{
			name: "100 elements",
			nums: randomSlice(100),
		},
		{
			name: "1000 elements",
			nums: randomSlice(1000),
		},
		{
			name: "10000 elements",
			nums: randomSlice(10000),
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CocktailSort(tc.nums)
			}
		})
	}
}

func BenchmarkCocktailSortBackSorted(b *testing.B) {
	testCases := []struct {
		name string
		nums []int
	}{
		{
			name: "10 elements",
			nums: backSortedSlice(10),
		},
		{
			name: "100 elements",
			nums: backSortedSlice(100),
		},
		{
			name: "1000 elements",
			nums: backSortedSlice(1000),
		},
		{
			name: "10000 elements",
			nums: backSortedSlice(10000),
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CocktailSort(tc.nums)
			}
		})
	}
}

func BenchmarkCocktailSortSorted(b *testing.B) {
	testCases := []struct {
		name string
		nums []int
	}{
		{
			name: "10 elements",
			nums: sortedSlice(10),
		},
		{
			name: "100 elements",
			nums: sortedSlice(100),
		},
		{
			name: "1000 elements",
			nums: sortedSlice(1000),
		},
		{
			name: "10000 elements",
			nums: sortedSlice(10000),
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CocktailSort(tc.nums)
			}
		})
	}
}

func randomSlice(sliceLen int) []int {
	res := make([]int, sliceLen)
	for i := 0; i < sliceLen; i++ {
		res[i] = rand.Intn(51)
	}
	return res
}

func backSortedSlice(sliceLen int) []int {
	res := make([]int, sliceLen)
	for i := sliceLen - 1; i >= 0; i-- {
		res[i] = i
	}
	return res
}

func sortedSlice(sliceLen int) []int {
	res := make([]int, sliceLen)
	for i := 0; i < sliceLen; i++ {
		res[i] = i
	}
	return res
}
