package edge

import "time"

type Edge struct {
	ID         string         `json:"id"`
	GraphID    string         `json:"graph_id"`
	Type       string         `json:"type"`
	FromID     string         `json:"from_id"`
	ToID       string         `json:"to_id"`
	Properties map[string]any `json:"properties"`
	Version    int64          `json:"version"`
	CreatedAt  time.Time      `json:"created_at"`
	Deleted    bool           `json:"deleted"`
}
