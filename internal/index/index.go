package index

import (
	"context"
	"fmt"
	"github.com/example/distributed-property-graph/internal/vertex"
	"strings"
	"sync"
)

type PropertyIndex struct {
	mu     sync.RWMutex
	values map[string]map[string]map[string]struct{}
}

func NewPropertyIndex() *PropertyIndex {
	return &PropertyIndex{values: map[string]map[string]map[string]struct{}{}}
}
func (i *PropertyIndex) Put(_ context.Context, v vertex.Vertex) error {
	for key, value := range v.Properties {
		token := strings.ToLower(v.Type + ":" + key + ":" + toString(value))
		if i.values[token] == nil {
			i.values[token] = map[string]map[string]struct{}{}
		}
		if i.values[token][v.GraphID] == nil {
			i.values[token][v.GraphID] = map[string]struct{}{}
		}
		i.values[token][v.GraphID][v.ID] = struct{}{}
	}
	return nil
}
func (i *PropertyIndex) Lookup(_ context.Context, g, typ, key string, value any) []string {
	ids := i.values[strings.ToLower(typ+":"+key+":"+toString(value))][g]
	out := []string{}
	for id := range ids {
		out = append(out, id)
	}
	return out
}
func toString(v any) string { return fmt.Sprintf("%v", v) }
