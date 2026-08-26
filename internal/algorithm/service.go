package algorithm

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/platform"
	"github.com/example/distributed-property-graph/internal/vertex"
	"sort"
	"sync"
)

type Service struct {
	mu       sync.RWMutex
	tasks    map[string]Task
	queue    chan string
	vertices vertex.Repository
	edges    edge.Repository
	clock    platform.Clock
}

func NewService(v vertex.Repository, e edge.Repository, clock platform.Clock) *Service {
	return &Service{tasks: map[string]Task{}, queue: make(chan string, 256), vertices: v, edges: e, clock: clock}
}
func (s *Service) Submit(ctx context.Context, g, name string, params map[string]any) (Task, error) {
	switch name {
	case "bfs", "connected_components", "pagerank", "triangle_count":
	default:
		return Task{}, platform.ErrInvalid
	}
	t := Task{ID: platform.NewID("algo"), GraphID: g, Name: name, Parameters: params, Status: StatusQueued, CreatedAt: s.clock.Now()}
	s.mu.Lock()
	s.tasks[t.ID] = t
	s.mu.Unlock()
	select {
	case s.queue <- t.ID:
		return t, nil
	case <-ctx.Done():
		return Task{}, ctx.Err()
	}
}
func (s *Service) Get(_ context.Context, id string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, platform.ErrNotFound
	}
	return t, nil
}
func (s *Service) Cancel(_ context.Context, id string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, platform.ErrNotFound
	}
	t.Cancel = true
	if t.Status == StatusQueued {
		t.Status = StatusCancelled
	}
	s.tasks[id] = t
	return t, nil
}
func (s *Service) Run(ctx context.Context, workers int) {
	if workers < 1 {
		workers = 1
	}
	for i := 0; i < workers; i++ {
		go s.worker(ctx)
	}
}
func (s *Service) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case id := <-s.queue:
			s.execute(ctx, id)
		}
	}
}
func (s *Service) execute(ctx context.Context, id string) {
	s.mu.Lock()
	t := s.tasks[id]
	if t.Cancel {
		s.mu.Unlock()
		return
	}
	t.Status = StatusRunning
	s.tasks[id] = t
	s.mu.Unlock()
	var result any
	var err error
	switch t.Name {
	case "bfs":
		result, err = s.bfs(ctx, t)
	case "connected_components":
		result, err = s.components(ctx, t)
	case "pagerank":
		result, err = s.pageRank(ctx, t, 20, 0.85)
	case "triangle_count":
		result, err = s.triangles(ctx, t)
	}
	s.mu.Lock()
	t = s.tasks[id]
	now := s.clock.Now()
	t.FinishedAt = &now
	if t.Cancel {
		t.Status = StatusCancelled
	} else if err != nil {
		t.Status = StatusFailed
		t.Error = err.Error()
	} else {
		t.Status = StatusCompleted
		t.Result = result
	}
	s.tasks[id] = t
	s.mu.Unlock()
}
func (s *Service) bfs(ctx context.Context, t Task) (any, error) {
	start, _ := t.Parameters["start"].(string)
	if start == "" {
		return nil, platform.ErrInvalid
	}
	queue := []string{start}
	seen := map[string]bool{start: true}
	order := []string{}
	for len(queue) > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)
		for _, e := range s.edges.ListFrom(ctx, t.GraphID, current) {
			if !seen[e.ToID] {
				seen[e.ToID] = true
				queue = append(queue, e.ToID)
			}
		}
	}
	return map[string]any{"order": order}, nil
}
func (s *Service) components(ctx context.Context, t Task) (any, error) {
	vertices := s.vertices.List(ctx, t.GraphID)
	component := map[string]int{}
	number := 0
	for _, v := range vertices {
		if _, ok := component[v.ID]; ok {
			continue
		}
		number++
		queue := []string{v.ID}
		component[v.ID] = number
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			for _, e := range s.edges.ListFrom(ctx, t.GraphID, current) {
				if _, ok := component[e.ToID]; !ok {
					component[e.ToID] = number
					queue = append(queue, e.ToID)
				}
			}
		}
	}
	return map[string]any{"components": component, "count": number}, nil
}
func (s *Service) pageRank(ctx context.Context, t Task, iterations int, damping float64) (any, error) {
	vertices := s.vertices.List(ctx, t.GraphID)
	if len(vertices) == 0 {
		return map[string]float64{}, nil
	}
	rank := map[string]float64{}
	base := 1 / float64(len(vertices))
	for _, v := range vertices {
		rank[v.ID] = base
	}
	for iteration := 0; iteration < iterations; iteration++ {
		next := map[string]float64{}
		for _, v := range vertices {
			next[v.ID] = (1 - damping) * base
		}
		for _, v := range vertices {
			out := s.edges.ListFrom(ctx, t.GraphID, v.ID)
			if len(out) == 0 {
				continue
			}
			share := damping * rank[v.ID] / float64(len(out))
			for _, e := range out {
				next[e.ToID] += share
			}
		}
		rank = next
	}
	return rank, nil
}
func (s *Service) triangles(ctx context.Context, t Task) (any, error) {
	adj := map[string]map[string]bool{}
	for _, e := range s.edges.List(ctx, t.GraphID) {
		if adj[e.FromID] == nil {
			adj[e.FromID] = map[string]bool{}
		}
		adj[e.FromID][e.ToID] = true
	}
	count := 0
	keys := make([]string, 0, len(adj))
	for key := range adj {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, a := range keys {
		for b := range adj[a] {
			for c := range adj[b] {
				if adj[c][a] {
					count++
				}
			}
		}
	}
	return map[string]int{"triangles": count / 3}, nil
}

var _ = fmt.Sprintf
