package edge

import (
	"sort"
	"strings"
)

type Filter struct {
	Type   string
	FromID string
	ToID   string
	Limit  int
}

func Match(e Edge, f Filter) bool {
	if f.Type != "" && e.Type != f.Type {
		return false
	}
	if f.FromID != "" && e.FromID != f.FromID {
		return false
	}
	if f.ToID != "" && e.ToID != f.ToID {
		return false
	}
	return true
}
func FilterEdges(values []Edge, f Filter) []Edge {
	out := make([]Edge, 0, len(values))
	for _, e := range values {
		if Match(e, f) {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	if f.Limit > 0 && len(out) > f.Limit {
		out = out[:f.Limit]
	}
	return out
}
func NormalizeType(value string) string { return strings.ToLower(strings.TrimSpace(value)) }
