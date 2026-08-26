package edge

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"sync"
)

type Repository interface {
	Put(context.Context, Edge) error
	Get(context.Context, string, string) (Edge, error)
	ListFrom(context.Context, string, string) []Edge
	List(context.Context, string) []Edge
}
type MemoryRepository struct {
	mu     sync.RWMutex
	values map[string]Edge
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{values: map[string]Edge{}} }
func (r *MemoryRepository) Put(_ context.Context, e Edge) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if old, ok := r.values[e.GraphID+"/"+e.ID]; ok && old.Version >= e.Version {
		return platform.ErrConflict
	}
	r.values[e.GraphID+"/"+e.ID] = e
	return nil
}
func (r *MemoryRepository) Get(_ context.Context, g, id string) (Edge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.values[g+"/"+id]
	if !ok || e.Deleted {
		return Edge{}, fmt.Errorf("edge %s: %w", id, platform.ErrNotFound)
	}
	return e, nil
}
func (r *MemoryRepository) ListFrom(_ context.Context, g, from string) []Edge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []Edge{}
	for _, e := range r.values {
		if e.GraphID == g && e.FromID == from && !e.Deleted {
			out = append(out, e)
		}
	}
	return out
}
func (r *MemoryRepository) List(_ context.Context, g string) []Edge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []Edge{}
	for _, e := range r.values {
		if e.GraphID == g && !e.Deleted {
			out = append(out, e)
		}
	}
	return out
}
