package vertex

import (
	"context"
	"errors"
	"testing"

	"github.com/example/distributed-property-graph/internal/platform"
)

func TestVertexNotFoundErrorWrapsSentinel(t *testing.T) {
	_, err := NewMemoryRepository().Get(context.Background(), "g", "missing")
	if !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestVertexDeleteNotFoundErrorWrapsSentinel(t *testing.T) {
	err := NewMemoryRepository().Delete(context.Background(), "g", "missing")
	if !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
