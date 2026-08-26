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
	out := map[string]any{}
	conflicts := []MergeConflict{}
	for k, v := range existing {
		out[k] = v
	}
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
