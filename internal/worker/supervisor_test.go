package worker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

type delayWorker struct{ done *atomic.Int32 }

func (w delayWorker) Run(context.Context) {
	time.Sleep(30 * time.Millisecond)
	w.done.Add(1)
}

func TestSupervisorRunAndWaitWaitsForWorkers(t *testing.T) {
	var done atomic.Int32
	s := &Supervisor{}
	s.Add(delayWorker{done: &done})
	s.RunAndWait(context.Background())
	if done.Load() != 1 {
		t.Fatal("RunAndWait returned before worker completed")
	}
}
