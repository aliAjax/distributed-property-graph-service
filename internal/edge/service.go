package edge

import (
	"context"
	"github.com/example/distributed-property-graph/internal/platform"
	"strings"
)

type VertexExists interface {
	Get(context.Context, string, string) (any, error)
}
type Service struct {
	repo  Repository
	clock platform.Clock
}

func NewService(repo Repository, clock platform.Clock) *Service {
	return &Service{repo: repo, clock: clock}
}
func (s *Service) Upsert(ctx context.Context, g, id, typ, from, to string, props map[string]any, version int64) (Edge, error) {
	if g == "" || id == "" || strings.TrimSpace(typ) == "" || from == "" || to == "" {
		return Edge{}, platform.ErrInvalid
	}
	if version < 1 {
		version = 1
	}
	e := Edge{ID: id, GraphID: g, Type: typ, FromID: from, ToID: to, Properties: props, Version: version, CreatedAt: s.clock.Now()}
	if old, err := s.repo.Get(ctx, g, id); err == nil {
		e.CreatedAt = old.CreatedAt
		if version <= old.Version {
			return Edge{}, platform.ErrConflict
		}
	}
	if err := s.repo.Put(ctx, e); err != nil {
		return Edge{}, err
	}
	return e, nil
}
