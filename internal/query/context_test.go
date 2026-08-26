package query

import (
	"context"
	"testing"
	"time"

	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/vertex"
)

func TestApplyLimitsPropagatesParentCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()
	child, stop, err := ApplyLimits(parent, Request{}, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	select {
	case <-child.Done():
	case <-time.After(50 * time.Millisecond):
		t.Fatal("child context did not inherit parent cancellation")
	}
}

func TestTraverseHonorsCancellation(t *testing.T) {
	v := vertex.NewMemoryRepository()
	e := edge.NewMemoryRepository()
	ctx := context.Background()
	_ = v.Put(ctx, vertex.Vertex{GraphID: "g", ID: "a", Version: 1})
	_ = e.Put(ctx, edge.Edge{GraphID: "g", ID: "e", FromID: "a", ToID: "b", Version: 1})
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	svc := NewService(v, e)
	_, err := svc.Traverse(cancelled, Request{GraphID: "g", StartVertex: "a", Depth: 10})
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}
