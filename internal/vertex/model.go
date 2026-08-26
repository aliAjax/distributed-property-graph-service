package vertex

import "time"

type Vertex struct {
	ID         string         `json:"id"`
	GraphID    string         `json:"graph_id"`
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`
	Version    int64          `json:"version"`
	CreatedAt  time.Time      `json:"created_at"`
	Deleted    bool           `json:"deleted"`
}

func (v Vertex) Property(name string) (any, bool) { value, ok := v.Properties[name]; return value, ok }

func CloneVertices(values []Vertex) []Vertex { return values }
