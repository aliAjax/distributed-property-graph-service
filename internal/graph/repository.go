package graph

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"sync"
)

type Repository interface {
	Save(context.Context, Graph) error
	Get(context.Context, string) (Graph, error)
	List(context.Context) ([]Graph, error)
}
type MemoryRepository struct {
	mu     sync.RWMutex
	values map[string]Graph
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{values: map[string]Graph{}} }
func (r *MemoryRepository) Save(_ context.Context, g Graph) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[g.ID] = g
	return nil
}
func (r *MemoryRepository) Get(_ context.Context, id string) (Graph, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.values[id]
	if !ok {
		return Graph{}, fmt.Errorf("graph %s: %w", id, platform.ErrNotFound)
	}
	return g, nil
}
func (r *MemoryRepository) List(_ context.Context) ([]Graph, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Graph, 0, len(r.values))
	for _, g := range r.values {
		out = append(out, g)
	}
	return out, nil
}
