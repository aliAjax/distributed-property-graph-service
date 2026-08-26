package algorithm

import (
	"sync"
	"testing"
	"time"
)

func TestPriorityQueueConcurrentEnqueueDequeue(t *testing.T) {
	q := &PriorityQueue{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 500; n++ {
			q.Enqueue("id", n)
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 500; n++ {
			_ = q.Dequeue()
		}
	}()
	close(start)
	wg.Wait()
}

func TestMetricsConcurrentCompleteAndSnapshot(t *testing.T) {
	m := &Metrics{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 500; n++ {
			m.Complete(time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 500; n++ {
			_ = m.Snapshot()
		}
	}()
	close(start)
	wg.Wait()
}

func TestMetricsConcurrentStartAndFail(t *testing.T) {
	m := &Metrics{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 500; n++ {
			m.Start()
			m.Fail()
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for n := 0; n < 500; n++ {
			m.Start()
			m.Fail()
		}
	}()
	close(start)
	wg.Wait()
}
