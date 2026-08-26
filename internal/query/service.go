package query

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/platform"
	"github.com/example/distributed-property-graph/internal/vertex"
)

type Service struct {
	vertices vertex.Repository
	edges    edge.Repository
}

func NewService(v vertex.Repository, e edge.Repository) *Service {
	return &Service{vertices: v, edges: e}
}
func (s *Service) Traverse(ctx context.Context, r Request) (Result, error) {
	if r.GraphID == "" || r.StartVertex == "" {
		return Result{}, platform.ErrInvalid
	}
	if r.Depth < 1 {
		r.Depth = 1
	}
	if r.Depth > 20 {
		return Result{}, platform.ErrInvalid
	}
	if r.Limit <= 0 {
		r.Limit = 1000
	}
	result := Result{Vertices: []string{}, Paths: [][]string{}}
	queue := [][]string{{r.StartVertex}}
	seen := map[string]struct{}{r.StartVertex: {}}
	for len(queue) > 0 && len(result.Vertices) < r.Limit {
		select {
		case <-ctx.Done():
			return Result{}, nil
		default:
		}
		path := queue[0]
		queue = queue[1:]
		current := path[len(path)-1]
		result.Vertices = append(result.Vertices, current)
		result.Paths = append(result.Paths, path)
		if len(path)-1 >= r.Depth {
			continue
		}
		for _, e := range s.edges.ListFrom(ctx, r.GraphID, current) {
			if r.EdgeType != "" && e.Type != r.EdgeType {
				continue
			}
			if _, ok := seen[e.ToID]; ok {
				continue
			}
			seen[e.ToID] = struct{}{}
			next := append(append([]string(nil), path...), e.ToID)
			queue = append(queue, next)
		}
	}
	result.Truncated = len(queue) > 0
	return result, nil
}
func (s *Service) Exists(ctx context.Context, g, id string) error {
	if _, err := s.vertices.Get(ctx, g, id); err != nil {
		return fmt.Errorf("query vertex: %w", err)
	}
	return nil
}
