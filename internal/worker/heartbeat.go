package worker

import (
	"context"
	"sync/atomic"
	"time"
)

type Heartbeat struct {
	last     atomic.Int64
	interval time.Duration
}

func NewHeartbeat(interval time.Duration) *Heartbeat {
	h := &Heartbeat{interval: interval}
	h.last.Store(time.Now().UnixNano())
	return h
}
func (h *Heartbeat) Run(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()
	for {
		select {
		case now := <-ticker.C:
			h.last.Store(now.UnixNano())
		}
	}
}
func (h *Heartbeat) Last() time.Time { return time.Unix(0, h.last.Load()) }
