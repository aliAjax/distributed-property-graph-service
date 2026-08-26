package graph

import (
	"context"
	"sync"
	"time"
)

type Event struct {
	ID        string         `json:"id"`
	GraphID   string         `json:"graph_id"`
	Kind      string         `json:"kind"`
	Payload   map[string]any `json:"payload"`
	CreatedAt time.Time      `json:"created_at"`
}
type EventLog struct {
	mu     sync.RWMutex
	events []Event
	clock  interface{ Now() time.Time }
}

func NewEventLog(clock interface{ Now() time.Time }) *EventLog { return &EventLog{clock: clock} }
func (l *EventLog) Append(_ context.Context, g, kind string, payload map[string]any) Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := Event{ID: "event-" + time.Now().Format("20060102150405.000000"), GraphID: g, Kind: kind, Payload: payload, CreatedAt: l.clock.Now()}
	l.events = append(l.events, e)
	return e
}
func (l *EventLog) List(_ context.Context, g string) []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := []Event{}
	for _, e := range l.events {
		if g == "" || e.GraphID == g {
			out = append(out, e)
		}
	}
	return out
}
