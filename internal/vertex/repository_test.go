package vertex

import (
	"context"
	"testing"
)

func TestMemoryRepositoryPreservesProperties(t *testing.T) {
	repository := NewMemoryRepository()
	value := Vertex{ID: "a", GraphID: "g", Type: "service", Version: 1, Properties: map[string]any{"name": "api"}}
	if err := repository.Put(context.Background(), value); err != nil {
		t.Fatal(err)
	}
	stored, err := repository.Get(context.Background(), "g", "a")
	if err != nil {
		t.Fatal(err)
	}
	if stored.Properties["name"] != "api" {
		t.Fatalf("property was not preserved: %+v", stored.Properties)
	}
	stored.Properties["name"] = "changed"
	again, err := repository.Get(context.Background(), "g", "a")
	if err != nil {
		t.Fatal(err)
	}
	if again.Properties["name"] != "api" {
		t.Fatal("repository leaked mutable property map")
	}
}
