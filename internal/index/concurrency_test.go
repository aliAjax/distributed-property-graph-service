package index

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/shard"
	"github.com/example/distributed-property-graph/internal/vertex"
)

func TestPropertyIndexConcurrentReadWrite(t *testing.T) {
	idx := NewPropertyIndex()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 200; n++ {
			_ = idx.Put(context.Background(), vertex.Vertex{GraphID: "g", ID: "v", Type: "person", Properties: map[string]any{"name": n}})
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 200; n++ {
			_ = idx.Lookup(context.Background(), "g", "person", "name", n)
		}
	}()
	close(start)
	wg.Wait()
}

func TestAdjacencyConcurrentReadWrite(t *testing.T) {
	adj := NewAdjacency()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 200; n++ {
			adj.Put(context.Background(), edgeEdge(n))
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 200; n++ {
			_ = adj.From(context.Background(), "g", "a")
		}
	}()
	close(start)
	wg.Wait()
}

func TestPropertyIndexConcurrentPut(t *testing.T) {
	idx := NewPropertyIndex()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	for n := 0; n < 2; n++ {
		go func(offset int) {
			defer wg.Done()
			<-start
			for i := 0; i < 200; i++ {
				_ = idx.Put(context.Background(), vertex.Vertex{GraphID: "g", ID: fmt.Sprintf("v%d-%d", offset, i), Type: "person", Properties: map[string]any{"name": i}})
			}
		}(n)
	}
	close(start)
	wg.Wait()
}

func TestAdjacencyConcurrentPut(t *testing.T) {
	adj := NewAdjacency()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	for n := 0; n < 2; n++ {
		go func(offset int) {
			defer wg.Done()
			<-start
			for i := 0; i < 200; i++ {
				adj.Put(context.Background(), edge.Edge{GraphID: "g", FromID: "a", ToID: "b", ID: fmt.Sprintf("e%d-%d", offset, i)})
			}
		}(n)
	}
	close(start)
	wg.Wait()
}

func TestCombinedIndexAdjacencyAndRingRace(t *testing.T) {
	idx := NewPropertyIndex()
	adj := NewAdjacency()
	ring := shard.NewRing(8)
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 200; n++ {
			_ = idx.Put(context.Background(), vertex.Vertex{GraphID: "g", ID: fmt.Sprintf("v%d", n), Type: "person", Properties: map[string]any{"name": n}})
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 200; n++ {
			adj.Put(context.Background(), edge.Edge{GraphID: "g", FromID: "a", ToID: "b", ID: fmt.Sprintf("e%d", n)})
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 200; n++ {
			ring.AddPoint(uint32(n))
			_ = ring.Lookup("shard-key")
		}
	}()
	close(start)
	wg.Wait()
}

func edgeEdge(n int) edge.Edge { return edge.Edge{GraphID: "g", FromID: "a", ToID: "b", ID: fmt.Sprintf("e%d", n)} }
