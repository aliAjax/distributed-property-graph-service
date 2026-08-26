package query

import (
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"strings"
)

func Validate(r Request) error {
	if strings.TrimSpace(r.GraphID) == "" || strings.TrimSpace(r.StartVertex) == "" {
		return platform.ErrInvalid
	}
	if r.Depth < 0 || r.Depth > 32 {
		return fmt.Errorf("depth: %w", platform.ErrInvalid)
	}
	if r.Limit < 0 || r.Limit > 100000 {
		return fmt.Errorf("limit: %w", platform.ErrInvalid)
	}
	return nil
}
func Normalize(r Request) Request {
	if r.Depth == 0 {
		r.Depth = 1
	}
	if r.Limit == 0 {
		r.Limit = 1000
	}
	r.EdgeType = strings.TrimSpace(r.EdgeType)
	return r
}
