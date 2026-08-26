package schema

type Property struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Indexed  bool   `json:"indexed"`
	Unique   bool   `json:"unique"`
}
type VertexType struct {
	Name       string     `json:"name"`
	Properties []Property `json:"properties"`
}
type EdgeType struct {
	Name       string     `json:"name"`
	From       []string   `json:"from"`
	To         []string   `json:"to"`
	Properties []Property `json:"properties"`
}
type Schema struct {
	GraphID   string       `json:"graph_id"`
	Version   int64        `json:"version"`
	Vertices  []VertexType `json:"vertices"`
	Edges     []EdgeType   `json:"edges"`
	Published bool         `json:"published"`
}
