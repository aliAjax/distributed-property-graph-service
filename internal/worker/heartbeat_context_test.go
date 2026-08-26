package worker

import (
	"context"
	"testing"
	"time"
)

func TestHeartbeatStopsOnCancel(t *testing.T) {
	h := NewHeartbeat(time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	go func() {
		h.Run(ctx)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("heartbeat did not stop after context cancellation")
	}
}
