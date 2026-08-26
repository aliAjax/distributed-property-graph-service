package graph

import (
	"context"
	"testing"

	"github.com/example/distributed-property-graph/internal/platform"
)

func TestPublishTransitionsToPublished(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, platform.SystemClock{})
	g, err := svc.Create(context.Background(), "g")
	if err != nil {
		t.Fatal(err)
	}
	published, err := svc.Publish(context.Background(), g.ID)
	if err != nil {
		t.Fatal(err)
	}
	if published.Status != StatusPublished {
		t.Fatalf("expected published status, got %q", published.Status)
	}
	got, err := repo.Get(context.Background(), g.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusPublished {
		t.Fatalf("repository not persisted: %q", got.Status)
	}
}
