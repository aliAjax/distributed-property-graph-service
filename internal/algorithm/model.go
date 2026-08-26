package algorithm

import "time"

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

type Task struct {
	ID         string         `json:"id"`
	GraphID    string         `json:"graph_id"`
	Name       string         `json:"name"`
	Parameters map[string]any `json:"parameters"`
	Status     Status         `json:"status"`
	Result     any            `json:"result,omitempty"`
	Error      string         `json:"error,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	FinishedAt *time.Time     `json:"finished_at,omitempty"`
	Cancel     bool           `json:"cancel_requested"`
}
