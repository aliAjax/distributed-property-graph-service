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

func (m *Metrics) Start() { m.Started++ }
func (m *Metrics) Complete(duration time.Duration) {
	m.Completed++
	m.TotalDuration += duration
}
func (m *Metrics) Fail()   { m.Failed++ }
func (m *Metrics) Cancel() { m.mu.Lock(); defer m.mu.Unlock(); m.Cancelled++ }

type MetricsSnapshot struct {
	Started, Completed, Failed, Cancelled uint64
	TotalDuration                         time.Duration
}

func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{Started: m.Started, Completed: m.Completed, Failed: m.Failed, Cancelled: m.Cancelled, TotalDuration: m.TotalDuration}
}
