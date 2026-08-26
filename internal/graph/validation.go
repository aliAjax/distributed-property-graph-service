package graph

import (
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"strings"
)

type Issue struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Validate(g Graph) []Issue {
	out := []Issue{}
	if strings.TrimSpace(g.Name) == "" {
		out = append(out, Issue{"name", "required", "name is required"})
	}
	if len(g.Name) > 200 {
		out = append(out, Issue{"name", "length", "name is too long"})
	}
	return out
}
func EnsureValid(g Graph) error {
	issues := Validate(g)
	if len(issues) == 0 {
		return nil
	}
	return fmt.Errorf("%s: %w", issues[0].Message, platform.ErrInvalid)
}
