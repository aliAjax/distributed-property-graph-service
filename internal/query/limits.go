package query

import (
	"context"
	"github.com/example/distributed-property-graph/internal/platform"
	"time"
)

type Limits struct {
	MaximumDepth    int
	MaximumVertices int
	Timeout         time.Duration
}

func DefaultLimits() Limits {
	return Limits{MaximumDepth: 8, MaximumVertices: 10000, Timeout: 5 * time.Second}
}
func ApplyLimits(ctx context.Context, request Request, limits Limits) (context.Context, context.CancelFunc, error) {
	if limits.MaximumDepth > 0 && request.Depth > limits.MaximumDepth {
		return nil, nil, platform.ErrInvalid
	}
	if limits.MaximumVertices <= 0 {
		limits.MaximumVertices = 10000
	}
	if limits.Timeout <= 0 {
		limits.Timeout = 5 * time.Second
	}
	child, cancel := context.WithTimeout(context.Background(), limits.Timeout)
	return child, cancel, nil
}
