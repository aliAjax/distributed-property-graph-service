package snapshot

import (
	"context"
	"github.com/example/distributed-property-graph/internal/platform"
	"sync"
	"time"
)

type Commit struct {
	ID        int64     `json:"id"`
	GraphID   string    `json:"graph_id"`
	Kind      string    `json:"kind"`
	EntityID  string    `json:"entity_id"`
	CreatedAt time.Time `json:"created_at"`
}
type CommitLog struct {
	mu    sync.RWMutex
	next  int64
	items []Commit
	clock platform.Clock
}

func NewCommitLog(clock platform.Clock) *CommitLog { return &CommitLog{clock: clock} }
func (l *CommitLog) Append(_ context.Context, g, kind, id string) Commit {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.next++
	c := Commit{ID: l.next, GraphID: g, Kind: kind, EntityID: id, CreatedAt: l.clock.Now()}
	l.items = append(l.items, c)
	return c
}
func (l *CommitLog) Since(_ context.Context, graph string, id int64) []Commit {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := []Commit{}
	for _, c := range l.items {
		if c.GraphID == graph && c.ID > id {
			out = append(out, c)
		}
	}
	return out
}
