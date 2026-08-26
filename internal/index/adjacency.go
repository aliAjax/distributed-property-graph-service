package index

import (
	"context"
	"github.com/example/distributed-property-graph/internal/edge"
	"sync"
)

type Adjacency struct {
	mu  sync.RWMutex
	out map[string][]edge.Edge
}

func NewAdjacency() *Adjacency { return &Adjacency{out: map[string][]edge.Edge{}} }
func (a *Adjacency) Put(_ context.Context, e edge.Edge) {
	key := e.GraphID + "/" + e.FromID
	a.out[key] = append(a.out[key], e)
}
func (a *Adjacency) From(_ context.Context, g, id string) []edge.Edge {
	return a.out[g+"/"+id]
}
