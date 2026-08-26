package shard

import (
	"sync"
	"testing"
)

func TestRingConcurrentMutationAndLookup(t *testing.T) {
	r := NewRing(8)
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 500; n++ {
			r.AddPoint(uint32(n))
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 500; n++ {
			_ = r.Lookup("shard-key")
		}
	}()
	close(start)
	wg.Wait()
}
