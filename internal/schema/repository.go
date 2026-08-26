package schema

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"sync"
)

type Repository interface {
	Save(context.Context, Schema) error
	Get(context.Context, string) (Schema, error)
}
type MemoryRepository struct {
	mu     sync.RWMutex
	values map[string]Schema
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{values: map[string]Schema{}} }
func (r *MemoryRepository) Save(_ context.Context, s Schema) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[s.GraphID] = s
	return nil
}
func (r *MemoryRepository) Get(_ context.Context, id string) (Schema, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.values[id]
	if !ok {
		return Schema{}, fmt.Errorf("schema %s: %v", id, platform.ErrNotFound)
	}
	return s, nil
}
