package edge

import (
	"context"
	"errors"
	"testing"

	"github.com/example/distributed-property-graph/internal/platform"
)

func TestEdgeNotFoundErrorWrapsSentinel(t *testing.T) {
	_, err := NewMemoryRepository().Get(context.Background(), "g", "missing")
	if !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEdgePutConflictErrorWrapsSentinel(t *testing.T) {
	repo := NewMemoryRepository()
	e := Edge{GraphID: "g", ID: "e", Version: 1}
	if err := repo.Put(context.Background(), e); err != nil {
		t.Fatal(err)
	}
	if err := repo.Put(context.Background(), e); !errors.Is(err, platform.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
