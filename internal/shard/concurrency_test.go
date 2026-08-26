package shard

import (
	"context"
	"sync"
	"testing"

	"github.com/example/distributed-property-graph/internal/platform"
)

func TestShardConcurrentReadWrite(t *testing.T) {
	svc := NewService(8, platform.SystemClock{})
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 300; n++ {
			_ = svc.List(context.Background())
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 300; n++ {
			_, _ = svc.Move(context.Background(), n%8, "node-2")
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 300; n++ {
			_ = svc.Route("shard-key")
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 300; n++ {
			_ = svc.Drain(context.Background(), "node-1")
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 300; n++ {
			_, _ = svc.PlanRebalance(context.Background(), []string{"node-1", "node-2"})
		}
	}()
	close(start)
	wg.Wait()
}
