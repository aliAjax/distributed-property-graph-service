package snapshot

import (
	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/vertex"
	"sort"
)

type Diff struct {
	AddedVertices   []string `json:"added_vertices"`
	RemovedVertices []string `json:"removed_vertices"`
	AddedEdges      []string `json:"added_edges"`
	RemovedEdges    []string `json:"removed_edges"`
}

func DiffVertices(before, after []vertex.Vertex) Diff {
	b := map[string]bool{}
	a := map[string]bool{}
	for _, v := range before {
		b[v.ID] = true
	}
	for _, v := range after {
		a[v.ID] = true
	}
	d := Diff{}
	for id := range a {
		if !b[id] {
			d.AddedVertices = append(d.AddedVertices, id)
		}
	}
	for id := range b {
		if !a[id] {
			d.RemovedVertices = append(d.RemovedVertices, id)
		}
	}
	sort.Strings(d.AddedVertices)
	sort.Strings(d.RemovedVertices)
	return d
}
func DiffEdges(before, after []edge.Edge) Diff {
	b := map[string]bool{}
	a := map[string]bool{}
	for _, v := range before {
		b[v.ID] = true
	}
	for _, v := range after {
		a[v.ID] = true
	}
	d := Diff{}
	for id := range a {
		if !b[id] {
			d.AddedEdges = append(d.AddedEdges, id)
		}
	}
	for id := range b {
		if !a[id] {
			d.RemovedEdges = append(d.RemovedEdges, id)
		}
	}
	sort.Strings(d.AddedEdges)
	sort.Strings(d.RemovedEdges)
	return d
}
