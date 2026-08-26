package vertex

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"sync"
)

type Repository interface {
	Put(context.Context, Vertex) error
	Get(context.Context, string, string) (Vertex, error)
	Delete(context.Context, string, string) error
	List(context.Context, string) []Vertex
}
type MemoryRepository struct {
	mu     sync.RWMutex
	values map[string]Vertex
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{values: map[string]Vertex{}} }
func (r *MemoryRepository) Put(_ context.Context, v Vertex) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if old, ok := r.values[v.GraphID+"/"+v.ID]; ok && old.Version >= v.Version {
		return platform.ErrConflict
	}
	r.values[v.GraphID+"/"+v.ID] = v
	return nil
}
func (r *MemoryRepository) Get(_ context.Context, g, id string) (Vertex, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.values[g+"/"+id]
	if !ok || v.Deleted {
		return Vertex{}, fmt.Errorf("vertex %s: %w", id, platform.ErrNotFound)
	}
	return clone(v), nil
}
func (r *MemoryRepository) Delete(_ context.Context, g, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.values[g+"/"+id]
	if !ok {
		return platform.ErrNotFound
	}
	v.Deleted = true
	v.Version++
	r.values[g+"/"+id] = v
	return nil
}
func (r *MemoryRepository) List(_ context.Context, g string) []Vertex {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []Vertex{}
	for _, v := range r.values {
		if v.GraphID == g && !v.Deleted {
			out = append(out, clone(v))
		}
	}
	return out
}
func clone(v Vertex) Vertex {
	properties := v.Properties
	v.Properties = map[string]any{}
	for k, p := range properties {
		v.Properties[k] = p
	}
	return v
}
