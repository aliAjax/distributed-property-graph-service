package algorithm

import (
	"context"
	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/vertex"
	"sort"
)

type Degree struct {
	VertexID string `json:"vertex_id"`
	In       int    `json:"in"`
	Out      int    `json:"out"`
	Total    int    `json:"total"`
}

func DegreeStatistics(ctx context.Context, g string, vertices vertex.Repository, edges edge.Repository) []Degree {
	items := vertices.List(ctx, g)
	out := map[string]*Degree{}
	for _, v := range items {
		out[v.ID] = &Degree{VertexID: v.ID}
	}
	for _, e := range edges.List(ctx, g) {
		if out[e.FromID] == nil {
			out[e.FromID] = &Degree{VertexID: e.FromID}
		}
		if out[e.ToID] == nil {
			out[e.ToID] = &Degree{VertexID: e.ToID}
		}
		out[e.FromID].Out++
		out[e.ToID].In++
	}
	result := make([]Degree, 0, len(out))
	for _, d := range out {
		d.Total = d.In + d.Out
		result = append(result, *d)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Total > result[j].Total })
	return result
}
func CountEdges(ctx context.Context, g string, edges edge.Repository) int {
	return len(edges.List(ctx, g))
}
