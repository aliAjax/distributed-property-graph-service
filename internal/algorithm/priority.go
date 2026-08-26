package algorithm

import (
	"container/heap"
	"sync"
)

type Item struct {
	ID       string
	Priority int
	index    int
}
type PriorityQueue struct {
	mu    sync.Mutex
	items []*Item
}

// Len, Less, Swap, Push and Pop implement container/heap.Interface.
// They are only called while the caller of Enqueue/Dequeue holds p.mu
// (container/heap never retains the receiver lock itself), so they
// touch p.items without taking the lock again.
func (p *PriorityQueue) Len() int           { return len(p.items) }
func (p *PriorityQueue) Less(i, j int) bool { return p.items[i].Priority > p.items[j].Priority }
func (p *PriorityQueue) Swap(i, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
	p.items[i].index = i
	p.items[j].index = j
}
func (p *PriorityQueue) Push(value any) {
	item := value.(*Item)
	item.index = len(p.items)
	p.items = append(p.items, item)
}
func (p *PriorityQueue) Pop() any {
	old := p.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // allow the GC to reclaim the popped item
	p.items = old[:n-1]
	return item
}

// Enqueue inserts an item with the given priority. Safe for concurrent use.
func (p *PriorityQueue) Enqueue(id string, priority int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	heap.Push(p, &Item{ID: id, Priority: priority})
}

// Dequeue removes and returns the highest-priority item, or nil if empty.
// Safe for concurrent use.
func (p *PriorityQueue) Dequeue() *Item {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.items) == 0 {
		return nil
	}
	return heap.Pop(p).(*Item)
}
