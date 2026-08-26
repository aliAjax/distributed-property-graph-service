package vertex

import (
	"fmt"
	"sort"
	"strings"
)

type Filter struct {
	Type     string
	Property string
	Value    any
	Limit    int
}

func Match(v Vertex, f Filter) bool {
	if f.Type != "" && v.Type != f.Type {
		return false
	}
	if f.Property != "" {
		value, ok := v.Properties[f.Property]
		if !ok || strings.TrimSpace(toString(value)) != strings.TrimSpace(toString(f.Value)) {
			return false
		}
	}
	return true
}
func FilterVertices(values []Vertex, f Filter) []Vertex {
	out := make([]Vertex, 0, len(values))
	for _, v := range values {
		if Match(v, f) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	if f.Limit > 0 && len(out) > f.Limit {
		out = out[:f.Limit]
	}
	return out
}
func toString(v any) string { return fmt.Sprintf("%v", v) }
