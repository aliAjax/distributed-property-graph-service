package snapshot

import "time"

type Snapshot struct {
	ID        string    `json:"id"`
	GraphID   string    `json:"graph_id"`
	CommitID  int64     `json:"commit_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Status    string    `json:"status"`
}
