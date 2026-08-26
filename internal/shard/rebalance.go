package shard

import (
	"context"
	"github.com/example/distributed-property-graph/internal/platform"
	"sort"
)

type RebalancePlan struct {
	Moves  []Move  `json:"moves"`
	Before []Shard `json:"before"`
	After  []Shard `json:"after"`
}
type Move struct {
	ShardID int    `json:"shard_id"`
	From    string `json:"from"`
	To      string `json:"to"`
}

func (s *Service) PlanRebalance(_ context.Context, nodes []string) (RebalancePlan, error) {
	if len(nodes) == 0 {
		return RebalancePlan{}, platform.ErrInvalid
	}
	before := make([]Shard, 0, len(s.shards))
	for _, v := range s.shards {
		before = append(before, v)
	}
	sort.Slice(before, func(i, j int) bool { return before[i].ID < before[j].ID })
	after := append([]Shard(nil), before...)
	moves := []Move{}
	for i := range after {
		target := nodes[i%len(nodes)]
		if after[i].NodeID != target {
			moves = append(moves, Move{ShardID: after[i].ID, From: after[i].NodeID, To: target})
			after[i].NodeID = target
			after[i].Epoch++
		}
	}
	return RebalancePlan{Moves: moves, Before: before, After: after}, nil
}
