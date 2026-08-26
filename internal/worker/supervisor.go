package worker

import (
	"context"
	"sync"
)

type Runnable interface{ Run(context.Context) }
type Supervisor struct {
	mu      sync.Mutex
	workers []Runnable
}

func (s *Supervisor) Add(worker Runnable) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workers = append(s.workers, worker)
}
func (s *Supervisor) Run(ctx context.Context) {
	s.mu.Lock()
	workers := append([]Runnable(nil), s.workers...)
	s.mu.Unlock()
	for _, worker := range workers {
		go worker.Run(ctx)
	}
}
func (s *Supervisor) RunAndWait(ctx context.Context) {
	s.mu.Lock()
	workers := append([]Runnable(nil), s.workers...)
	s.mu.Unlock()
	var wg sync.WaitGroup
	for _, worker := range workers {
		go func(w Runnable) {
			defer wg.Done()
			w.Run(ctx)
		}(worker)
	}
	wg.Wait()
}
func (s *Supervisor) Count() int { s.mu.Lock(); defer s.mu.Unlock(); return len(s.workers) }
