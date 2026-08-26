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

func CloneEdge(e Edge) Edge {
	if e.Properties != nil {
		props := make(map[string]any, len(e.Properties))
		for k, p := range e.Properties {
			props[k] = p
		}
		e.Properties = props
	}
	return e
}

func CloneEdges(values []Edge) []Edge {
	if values == nil {
		return nil
	}
	out := make([]Edge, len(values))
	for i, e := range values {
		out[i] = CloneEdge(e)
	}
	return out
}
