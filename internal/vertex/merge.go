package vertex

import (
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
)

type MergeConflict struct {
	Field    string `json:"field"`
	Existing any    `json:"existing"`
	Incoming any    `json:"incoming"`
}

func MergeProperties(existing, incoming map[string]any) (map[string]any, []MergeConflict) {
	// Build a fresh map so callers never receive the same backing map they
	// passed in, and so a nil `existing` (e.g. a vertex/edge written without
	// properties) does not panic on write. Both inputs may be nil: ranging over
	// a nil map yields no entries and is safe.
	out := make(map[string]any, len(existing)+len(incoming))
	for k, v := range existing {
		out[k] = v
	}
	conflicts := []MergeConflict{}
	for k, v := range incoming {
		if old, ok := out[k]; ok && fmt.Sprintf("%v", old) != fmt.Sprintf("%v", v) {
			conflicts = append(conflicts, MergeConflict{Field: k, Existing: old, Incoming: v})
		} else {
			out[k] = v
		}
	}
	return out, conflicts
}
func EnsureNoConflicts(conflicts []MergeConflict) error {
	if len(conflicts) > 0 {
		return platform.ErrConflict
	}
	return nil
}
