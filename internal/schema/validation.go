package schema

import (
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"strings"
)

func ValidateProperty(p Property) error {
	if strings.TrimSpace(p.Name) == "" || p.Type == "" {
		return platform.ErrInvalid
	}
	switch p.Type {
	case "string", "integer", "number", "boolean", "timestamp", "json":
		return nil
	default:
		return fmt.Errorf("property type %s: %w", p.Type, platform.ErrInvalid)
	}
}
func ValidateSchema(s Schema) error {
	if s.GraphID == "" || len(s.Vertices) == 0 {
		return platform.ErrInvalid
	}
	names := map[string]bool{}
	for _, v := range s.Vertices {
		if names[v.Name] {
			return platform.ErrConflict
		}
		names[v.Name] = true
		for _, p := range v.Properties {
			if err := ValidateProperty(p); err != nil {
				return err
			}
		}
	}
	return nil
}
func Merge(base, overlay Schema) Schema {
	out := base
	for _, v := range overlay.Vertices {
		found := false
		for i := range out.Vertices {
			if out.Vertices[i].Name == v.Name {
				out.Vertices[i] = v
				found = true
			}
		}
		if !found {
			out.Vertices = append(out.Vertices, v)
		}
	}
	for _, e := range overlay.Edges {
		found := false
		for i := range out.Edges {
			if out.Edges[i].Name == e.Name {
				out.Edges[i] = e
				found = true
			}
		}
		if !found {
			out.Edges = append(out.Edges, e)
		}
	}
	out.Version++
	return out
}
