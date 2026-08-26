package edge

import (
	"context"
	"testing"
)

func TestFilterEdgesDoesNotMutateInput(t *testing.T) {
	values := []Edge{
		{ID: "e2", Type: "b"},
		{ID: "e1", Type: "a"},
	}
	filtered := FilterEdges(values, Filter{Type: "a"})
	if len(filtered) != 1 || filtered[0].ID != "e1" {
		t.Fatalf("unexpected filtered: %#v", filtered)
	}
	if len(values) != 2 || values[0].ID != "e2" || values[1].ID != "e1" {
		t.Fatalf("input was mutated: %#v", values)
	}
}

type staticRepo struct{ edges []Edge }

func (r *staticRepo) Put(context.Context, Edge) error { return nil }
func (r *staticRepo) Get(context.Context, string, string) (Edge, error) { return Edge{}, nil }
func (r *staticRepo) ListFrom(_ context.Context, g, from string) []Edge {
	out := []Edge{}
	for _, e := range r.edges {
		if e.FromID == from {
			out = append(out, e)
		}
	}
	return out
}
func (r *staticRepo) List(context.Context, string) []Edge { return append([]Edge(nil), r.edges...) }

func TestSimplePathsReturnIndependentSlices(t *testing.T) {
	repo := &staticRepo{edges: []Edge{
		{GraphID: "g", FromID: "a", ToID: "b"},
		{GraphID: "g", FromID: "a", ToID: "c"},
		{GraphID: "g", FromID: "b", ToID: "d"},
		{GraphID: "g", FromID: "c", ToID: "d"},
	}}
	paths := SimplePaths(context.Background(), repo, "g", "a", "d", 3)
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	paths[0][0] = "mutated"
	if paths[1][0] != "a" {
		t.Fatalf("paths share storage: %#v", paths)
	}
}

func TestCloneEdgesDoesNotAlias(t *testing.T) {
	values := []Edge{{ID: "e1"}}
	cloned := CloneEdges(values)
	cloned[0].ID = "mutated"
	if values[0].ID != "e1" {
		t.Fatal("CloneEdges returned aliased storage")
	}
}

func TestClonePathsDoesNotAlias(t *testing.T) {
	values := [][]string{{"a", "b"}}
	cloned := ClonePaths(values)
	cloned[0][0] = "mutated"
	if values[0][0] != "a" {
		t.Fatal("ClonePaths returned aliased storage")
	}
}
