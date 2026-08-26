package edge

import (
	"context"
	"testing"

	"github.com/example/distributed-property-graph/internal/platform"
)

func TestEdgeUpsertInitializesNilProperties(t *testing.T) {
	svc := NewService(NewMemoryRepository(), platform.SystemClock{})
	e, err := svc.Upsert(context.Background(), "g", "e1", "knows", "a", "b", nil, 1)
	if err != nil {
		t.Fatal(err)
	}
	if e.Properties == nil {
		t.Fatal("expected non-nil properties")
	}
}
