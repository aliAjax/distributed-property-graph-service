package shard

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

type Ring struct {
	mu     sync.RWMutex
	points map[uint32]int
	sorted []uint32
}

func NewRing(count int) *Ring {
	r := &Ring{points: map[uint32]int{}}
	for i := 0; i < count; i++ {
		point := crc32.ChecksumIEEE([]byte(fmt.Sprintf("shard-%d", i)))
		r.points[point] = i
		r.sorted = append(r.sorted, point)
	}
	sort.Slice(r.sorted, func(i, j int) bool { return r.sorted[i] < r.sorted[j] })
	return r
}
func (r *Ring) Lookup(key string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.sorted) == 0 {
		return -1
	}
	point := crc32.ChecksumIEEE([]byte(key))
	index := sort.Search(len(r.sorted), func(i int) bool { return r.sorted[i] >= point })
	if index == len(r.sorted) {
		index = 0
	}
	return r.points[r.sorted[index]]
}
