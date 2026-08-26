package vertex

import (
	"context"
	"testing"

	"github.com/example/distributed-property-graph/internal/platform"
)

func TestUpsertInitializesNilProperties(t *testing.T) {
	svc := NewService(NewMemoryRepository(), platform.SystemClock{})
	v, err := svc.Upsert(context.Background(), "g", "v1", "person", nil, 1)
	if err != nil {
		t.Fatal(err)
	}
	if v.Properties == nil {
		t.Fatal("expected non-nil properties")
	}
}

func TestMergePropertiesDoesNotPanicOnNilExisting(t *testing.T) {
	defer func() {
		if recover() != nil {
			t.Fatal("MergeProperties panicked on nil existing map")
		}
	}()
	out, conflicts := MergeProperties(nil, map[string]any{"age": 30})
	if len(conflicts) != 0 || out["age"] != 30 {
		t.Fatalf("unexpected merge result: %#v conflicts=%v", out, conflicts)
	}
}
