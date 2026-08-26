package algorithm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/platform"
	"github.com/example/distributed-property-graph/internal/vertex"
)

func TestRetryFailedTaskResetsState(t *testing.T) {
	svc := NewService(vertex.NewMemoryRepository(), edge.NewMemoryRepository(), platform.SystemClock{})
	finished := time.Now()
	svc.tasks["t"] = Task{ID: "t", Status: StatusFailed, Error: "boom", Result: map[string]any{"x": 1}, FinishedAt: &finished}
	got, err := svc.Retry(context.Background(), "t")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusQueued || got.Error != "" || got.Result != nil || got.FinishedAt != nil {
		t.Fatalf("failed task was not reset: %#v", got)
	}
}

func TestRetryCompletedTaskRejected(t *testing.T) {
	svc := NewService(vertex.NewMemoryRepository(), edge.NewMemoryRepository(), platform.SystemClock{})
	svc.tasks["t"] = Task{ID: "t", Status: StatusCompleted}
	if _, err := svc.Retry(context.Background(), "t"); !errors.Is(err, platform.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
