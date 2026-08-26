package snapshot

import (
	"context"
	"sync"
	"time"
)

type CleanupStore struct {
	mu    sync.Mutex
	items map[string]Snapshot
}

func NewCleanupStore() *CleanupStore { return &CleanupStore{items: map[string]Snapshot{}} }
func (c *CleanupStore) Cleanup(_ context.Context, now time.Time) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for id, s := range c.items {
		if !now.Before(s.ExpiresAt) {
			delete(c.items, id)
			n++
		}
	}
	return n, nil
}
