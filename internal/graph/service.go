package graph

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
func (s *Service) Create(ctx context.Context, name string) (Graph, error) {
	if strings.TrimSpace(name) == "" {
		return Graph{}, platform.ErrInvalid
	}
	g := Graph{ID: platform.NewID("graph"), Name: name, Status: StatusDraft, Version: 1, CreatedAt: s.clock.Now()}
	if err := s.repo.Save(ctx, g); err != nil {
		return Graph{}, fmt.Errorf("save graph: %w", err)
	}
	return g, nil
}
func (s *Service) Publish(ctx context.Context, id string) (Graph, error) {
	g, err := s.repo.Get(ctx, id)
	if err != nil {
		return Graph{}, err
	}
	if g.Status != StatusDraft {
		return Graph{}, platform.ErrConflict
	}
	now := s.clock.Now()
	g.Status = StatusDraft
	g.PublishedAt = &now
	g.Version++
	return g, nil
}
func (s *Service) Get(ctx context.Context, id string) (Graph, error) { return s.repo.Get(ctx, id) }
func (s *Service) List(ctx context.Context) ([]Graph, error)         { return s.repo.List(ctx) }
