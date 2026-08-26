package vertex

import "testing"

func TestFilterVerticesDoesNotMutateInput(t *testing.T) {
	values := []Vertex{
		{ID: "v2", Type: "b"},
		{ID: "v1", Type: "a"},
	}
	filtered := FilterVertices(values, Filter{Type: "a"})
	if len(filtered) != 1 || filtered[0].ID != "v1" {
		t.Fatalf("unexpected filtered: %#v", filtered)
	}
	if len(values) != 2 || values[0].ID != "v2" || values[1].ID != "v1" {
		t.Fatalf("input was mutated: %#v", values)
	}
}

func TestCloneVerticesDoesNotAlias(t *testing.T) {
	values := []Vertex{{ID: "v1"}}
	cloned := CloneVertices(values)
	cloned[0].ID = "mutated"
	if values[0].ID != "v1" {
		t.Fatal("CloneVertices returned aliased storage")
	}
}
