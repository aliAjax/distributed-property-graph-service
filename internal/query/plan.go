package query

import (
	"fmt"
	"strings"
)

type Plan struct {
	Steps         []string `json:"steps"`
	EstimatedCost int64    `json:"estimated_cost"`
	Warnings      []string `json:"warnings"`
}

func BuildPlan(r Request) Plan {
	steps := []string{"lookup start vertex"}
	if r.EdgeType != "" {
		steps = append(steps, "filter edge type "+r.EdgeType)
	}
	steps = append(steps, fmt.Sprintf("traverse depth %d", r.Depth), fmt.Sprintf("limit %d", r.Limit))
	return Plan{Steps: steps, EstimatedCost: int64(r.Depth * r.Limit)}
}
func (p Plan) String() string { return strings.Join(p.Steps, " -> ") }
