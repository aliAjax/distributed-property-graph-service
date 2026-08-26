package snapshot

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"sync"
	"time"
)

type Service struct {
	mu     sync.RWMutex
	values map[string]Snapshot
	clock  platform.Clock
	commit int64
}

func NewService(clock platform.Clock) *Service {
	return &Service{values: map[string]Snapshot{}, clock: clock}
}
func (s *Service) Create(_ context.Context, g string, ttl time.Duration) (Snapshot, error) {
	if g == "" {
		return Snapshot{}, platform.ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commit++
	now := s.clock.Now()
	snap := Snapshot{ID: platform.NewID("snapshot"), GraphID: g, CommitID: s.commit, CreatedAt: now, ExpiresAt: now.Add(ttl), Status: "active"}
	s.values[snap.ID] = snap
	return snap, nil
}
func (s *Service) Get(_ context.Context, id string) (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.values[id]
	if !ok {
		return Snapshot{}, fmt.Errorf("snapshot %s: %w", id, platform.ErrNotFound)
	}
	if !time.Now().Before(v.ExpiresAt) {
		return Snapshot{}, platform.ErrTimeout
	}
	return v, nil
}
func (s *Service) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.values[id]; !ok {
		return fmt.Errorf("snapshot: %w", platform.ErrNotFound)
	}
	delete(s.values, id)
	return nil
}
