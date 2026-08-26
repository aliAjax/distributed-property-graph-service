package snapshot

import (
	"context"
	"time"
)

type Store interface {
	Create(context.Context, string, time.Duration) (Snapshot, error)
	Get(context.Context, string) (Snapshot, error)
	Delete(context.Context, string) error
}
type Cleaner interface {
	Cleanup(context.Context, time.Time) (int, error)
}
