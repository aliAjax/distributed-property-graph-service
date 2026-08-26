package snapshot

import (
	"context"
	"errors"
	"testing"

	"github.com/example/distributed-property-graph/internal/platform"
)

func TestSnapshotGetNotFoundErrorWrapsSentinel(t *testing.T) {
	svc := NewService(platform.SystemClock{})
	_, err := svc.Get(context.Background(), "missing")
	if !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSnapshotDeleteNotFoundErrorWrapsSentinel(t *testing.T) {
	svc := NewService(platform.SystemClock{})
	err := svc.Delete(context.Background(), "missing")
	if !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
