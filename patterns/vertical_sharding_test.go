package patterns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestShardingSetAndGet tests Set and Get... by setting some values and getting them.
func TestPutAndGet(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap[string, int](BUCKETS)

	truthMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}

	for k, v := range truthMap {
		sMap.Put(k, v)
	}

	for k, v := range truthMap {
		got, ok := sMap.Get(k)

		assert.Equal(t, v, got)
		assert.Equal(t, true, ok)
	}
}

// TestShardingKeys tests the Keys method by adding 5 values to the map and checking
// that each one exists in the keys list exactly once.
func TestKeys(t *testing.T) {
	expectedMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}
	expectedKeys := make([]string, 0, len(expectedMap))
	for k := range expectedMap {
		expectedKeys = append(expectedKeys, k)
	}

	sMap := NewShardedMap[string, int](17)
	for k, v := range expectedMap {
		sMap.Put(k, v)
	}
	actualKeys := sMap.Keys()

	assert.ElementsMatch(t, expectedKeys, actualKeys)
}

func TestValues(t *testing.T) {
	expectedMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}
	expectedValues := make([]int, 0, len(expectedMap))
	for _, v := range expectedMap {
		expectedValues = append(expectedValues, v)
	}

	sMap := NewShardedMap[string, int](17)
	for k, v := range expectedMap {
		sMap.Put(k, v)
	}
	actualKeys := sMap.Values()

	assert.ElementsMatch(t, expectedValues, actualKeys)
}

// TestShardingDelete tests the Delete method by adding and then removing five values.
func TestRemove(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap[string, int](BUCKETS)

	truthMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}

	for k, v := range truthMap {
		sMap.Put(k, v)
	}

	keys := sMap.Keys()
	for _, key := range keys {
		sMap.Remove(key)
	}

	assert.Equal(t, 0, len(sMap.Keys()))
}

func TestClear(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap[string, int](BUCKETS)

	truthMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}

	for k, v := range truthMap {
		sMap.Put(k, v)
	}

	sMap.Clear()

	assert.Equal(t, 0, len(sMap.Keys()))
}
