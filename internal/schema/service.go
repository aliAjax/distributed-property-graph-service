package schema

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"strings"
)

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) Publish(ctx context.Context, value Schema) (Schema, error) {
	if value.GraphID == "" || len(value.Vertices) == 0 {
		return Schema{}, platform.ErrInvalid
	}
	for _, v := range value.Vertices {
		if strings.TrimSpace(v.Name) == "" {
			return Schema{}, platform.ErrInvalid
		}
	}
	current, err := s.repo.Get(ctx, value.GraphID)
	if err == nil && current.Published {
		return Schema{}, platform.ErrConflict
	}
	value.Version++
	value.Published = true
	if err := s.repo.Save(ctx, value); err != nil {
		return Schema{}, fmt.Errorf("save schema: %v", err)
	}
	return value, nil
}
func (s *Service) Get(ctx context.Context, id string) (Schema, error) { return s.repo.Get(ctx, id) }
