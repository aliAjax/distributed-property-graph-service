package vertex

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"strings"
)

type Service struct {
	repo  Repository
	clock platform.Clock
}

func NewService(repo Repository, clock platform.Clock) *Service {
	return &Service{repo: repo, clock: clock}
}
func (s *Service) Upsert(ctx context.Context, g, id, typ string, props map[string]any, version int64) (Vertex, error) {
	if g == "" || strings.TrimSpace(id) == "" || typ == "" {
		return Vertex{}, platform.ErrInvalid
	}
	if version < 1 {
		version = 1
	}
	v := Vertex{ID: id, GraphID: g, Type: typ, Properties: props, Version: version, CreatedAt: s.clock.Now()}
	if old, err := s.repo.Get(ctx, g, id); err == nil {
		v.CreatedAt = old.CreatedAt
		if version <= old.Version {
			return Vertex{}, fmt.Errorf("version: %d: %w", version, platform.ErrConflict)
		}
	}
	if err := s.repo.Put(ctx, v); err != nil {
		return Vertex{}, err
	}
	return v, nil
}
func (s *Service) Get(ctx context.Context, g, id string) (Vertex, error) {
	return s.repo.Get(ctx, g, id)
}
func (s *Service) Delete(ctx context.Context, g, id string) error { return s.repo.Delete(ctx, g, id) }
