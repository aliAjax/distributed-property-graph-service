package shard

import (
	"context"
	"github.com/example/distributed-property-graph/internal/platform"
	"hash/fnv"
	"sync"
	"time"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusMoving   Status = "moving"
	StatusReadonly Status = "readonly"
)

type Shard struct {
	ID           int       `json:"id"`
	NodeID       string    `json:"node_id"`
	Status       Status    `json:"status"`
	Epoch        int64     `json:"epoch"`
	LeaseExpires time.Time `json:"lease_expires"`
}
type Service struct {
	mu     sync.RWMutex
	count  int
	shards map[int]Shard
	clock  platform.Clock
}

func NewService(count int, clock platform.Clock) *Service {
	s := &Service{count: count, shards: map[int]Shard{}, clock: clock}
	for i := 0; i < count; i++ {
		s.shards[i] = Shard{ID: i, NodeID: "node-1", Status: StatusActive, Epoch: 1, LeaseExpires: clock.Now().Add(time.Minute)}
	}
	return s
}
func (s *Service) Route(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32() % uint32(s.count))
}
func (s *Service) List(_ context.Context) []Shard {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Shard, 0, len(s.shards))
	for _, v := range s.shards {
		out = append(out, v)
	}
	return out
}
func (s *Service) Move(_ context.Context, id int, node string) (Shard, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.shards[id]
	if !ok {
		return Shard{}, platform.ErrNotFound
	}
	v.Status = StatusMoving
	v.NodeID = node
	v.Epoch++
	v.LeaseExpires = s.clock.Now().Add(time.Minute)
	v.Status = StatusActive
	s.shards[id] = v
	return v, nil
}
func (s *Service) Drain(_ context.Context, node string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for id, v := range s.shards {
		if v.NodeID == node {
			v.Status = StatusReadonly
			s.shards[id] = v
			n++
		}
	}
	return n
}
