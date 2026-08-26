package algorithm

import (
	"sync"
	"time"
)

type Metrics struct {
	mu                                    sync.RWMutex
	Started, Completed, Failed, Cancelled uint64
	TotalDuration                         time.Duration
}

func (m *Metrics) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Started++
}
func (m *Metrics) Complete(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Completed++
	m.TotalDuration += duration
}
func (m *Metrics) Fail() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Failed++
}
func (m *Metrics) Cancel() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Cancelled++
}

type MetricsSnapshot struct {
	Started, Completed, Failed, Cancelled uint64
	TotalDuration                         time.Duration
}

// Snapshot returns a consistent point-in-time copy of all counters.
func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return MetricsSnapshot{Started: m.Started, Completed: m.Completed, Failed: m.Failed, Cancelled: m.Cancelled, TotalDuration: m.TotalDuration}
}
