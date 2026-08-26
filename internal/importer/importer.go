package importer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/platform"
	"github.com/example/distributed-property-graph/internal/vertex"
	"io"
)

type Record struct {
	Kind       string         `json:"kind"`
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	FromID     string         `json:"from_id,omitempty"`
	ToID       string         `json:"to_id,omitempty"`
	Properties map[string]any `json:"properties"`
}
type Failure struct {
	Line  int    `json:"line"`
	Error string `json:"error"`
}
type Summary struct {
	Accepted int       `json:"accepted"`
	Rejected int       `json:"rejected"`
	LastLine int       `json:"last_line"`
	Failures []Failure `json:"failures"`
}
type Importer struct {
	vertices *vertex.Service
	edges    *edge.Service
	maximum  int
}

func New(v *vertex.Service, e *edge.Service, maximum int) *Importer {
	return &Importer{vertices: v, edges: e, maximum: maximum}
}
func (i *Importer) Import(ctx context.Context, g string, r io.Reader, startLine int) (Summary, error) {
	summary := Summary{Failures: []Failure{}}
	scanner := bufio.NewScanner(io.LimitReader(r, 128<<20))
	scanner.Buffer(make([]byte, 64<<10), 1<<20)
	line := 0
	for scanner.Scan() {
		line++
		if line <= startLine {
			continue
		}
		if summary.Accepted+summary.Rejected >= i.maximum {
			return summary, platform.ErrInvalid
		}
		var record Record
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			summary.Rejected++
			summary.Failures = append(summary.Failures, Failure{line, err.Error()})
			continue
		}
		var err error
		switch record.Kind {
		case "vertex":
			_, err = i.vertices.Upsert(ctx, g, record.ID, record.Type, record.Properties, 1)
		case "edge":
			_, err = i.edges.Upsert(ctx, g, record.ID, record.Type, record.FromID, record.ToID, record.Properties, 1)
		default:
			err = platform.ErrInvalid
		}
		if err != nil {
			summary.Rejected++
			summary.Failures = append(summary.Failures, Failure{line, err.Error()})
		} else {
			summary.Accepted++
		}
		summary.LastLine = line
	}
	if err := scanner.Err(); err != nil {
		return summary, fmt.Errorf("scan import: %w", err)
	}
	return summary, nil
}
