package importer

import (
	"context"
	"strings"
	"testing"

	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/platform"
	"github.com/example/distributed-property-graph/internal/vertex"
)

func TestImportHonorsCancellation(t *testing.T) {
	v := vertex.NewService(vertex.NewMemoryRepository(), platform.SystemClock{})
	e := edge.NewService(edge.NewMemoryRepository(), platform.SystemClock{})
	imp := New(v, e, 1000)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := imp.Import(ctx, "g", strings.NewReader("{\"kind\":\"vertex\",\"id\":\"v1\",\"type\":\"person\"}\n"), 0)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}
