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

func CloneVertex(v Vertex) Vertex {
	if v.Properties != nil {
		props := make(map[string]any, len(v.Properties))
		for k, p := range v.Properties {
			props[k] = p
		}
		v.Properties = props
	}
	return v
}

func CloneVertices(values []Vertex) []Vertex {
	if values == nil {
		return nil
	}
	out := make([]Vertex, len(values))
	for i, v := range values {
		out[i] = CloneVertex(v)
	}
	return out
}
