package patterns

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"sync"

	"github.com/emirpasic/gods/v2/maps"
)

var _ maps.Map[string, int] = shardedMap[string, int]{}

type shard[K comparable, V any] struct {
	sync.RWMutex         // Compose from sync.RWMutex
	items        map[K]V // m contains the shard's data
}

type shardedMap[K comparable, V any] []*shard[K, V]

func NewShardedMap[K comparable, V any](nshards int) maps.Map[K, V] {
	shards := make([]*shard[K, V], nshards) // Initialize a *Shards slice

	for i := 0; i < nshards; i++ {
		shrd := make(map[K]V)
		shards[i] = &shard[K, V]{items: shrd} // A ShardedMap IS a slice!
	}

	return shardedMap[K, V](shards)
}

func (m shardedMap[K, V]) Remove(key K) {
	shrd := m.getShard(key)
	shrd.Lock()
	defer shrd.Unlock()

	delete(shrd.items, key)
}

func (m shardedMap[K, V]) Get(key K) (V, bool) {
	shrd := m.getShard(key)
	shrd.RLock()
	defer shrd.RUnlock()
	item, ok := shrd.items[key]
	return item, ok
}

func (m shardedMap[K, V]) Put(key K, value V) {
	shrd := m.getShard(key)
	shrd.Lock()
	defer shrd.Unlock()

	shrd.items[key] = value
}

func (m shardedMap[K, V]) Keys() []K {
	keys, _ := m.keysValues()
	return keys
}

func (m shardedMap[K, V]) Values() []V {
	_, values := m.keysValues()
	return values
}

func (m shardedMap[K, V]) Clear() {
	var wg sync.WaitGroup
	wg.Add(len(m))
	for _, shrd := range m {
		go func(s *shard[K, V]) {
			defer wg.Done()
			s.Lock()
			defer s.Unlock()
			s.items = make(map[K]V)
		}(shrd)
	}
	wg.Wait()
}

func (m shardedMap[K, V]) Empty() bool {
	return len(m.Keys()) == 0
}

func (m shardedMap[K, V]) Size() int {
	return len(m.Keys())
}

func (m shardedMap[K, V]) String() string {
	panic("not implemented")
}

func (m shardedMap[K, V]) keysValues() ([]K, []V) {
	var wg sync.WaitGroup
	wg.Add(len(m))
	ch := make(chan struct {
		ks []K
		vs []V
	}, len(m))
	for _, shrd := range m {
		go func(s *shard[K, V]) {
			s.RLock()
			keys := make([]K, 0, len(s.items))
			values := make([]V, 0, len(s.items))
			for key, value := range s.items {
				keys = append(keys, key)
				values = append(values, value)
			}
			s.RUnlock()

			ch <- struct {
				ks []K
				vs []V
			}{
				keys,
				values,
			}
			wg.Done()
		}(shrd)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var keys []K
	var values []V
	for kv := range ch {
		keys = append(keys, kv.ks...)
		values = append(values, kv.vs...)
	}
	return keys, values
}

func (m shardedMap[K, V]) getShardIndex(key K) int {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(key)
	hash := fnv.New32a()     // Get a hash implementation from "hash/fnv"
	hash.Write(buf.Bytes())  // Write bytes to the hash
	sum := int(hash.Sum32()) // Get the resulting checksum
	return sum % len(m)      // Mod by len(m) to get index
}

func (m shardedMap[K, V]) getShard(key K) *shard[K, V] {
	index := m.getShardIndex(key)
	return m[index]
}
