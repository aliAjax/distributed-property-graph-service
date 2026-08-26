package query

type Request struct {
	GraphID       string `json:"graph_id"`
	StartVertex   string `json:"start_vertex"`
	Depth         int    `json:"depth"`
	EdgeType      string `json:"edge_type,omitempty"`
	PropertyKey   string `json:"property_key,omitempty"`
	PropertyValue any    `json:"property_value,omitempty"`
	Limit         int    `json:"limit"`
}
type Result struct {
	Vertices   []string   `json:"vertices"`
	Paths      [][]string `json:"paths"`
	Truncated  bool       `json:"truncated"`
	SnapshotID string     `json:"snapshot_id,omitempty"`
}
