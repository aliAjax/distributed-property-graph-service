package platform

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Cursor struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Shard     int       `json:"shard"`
}

func EncodeCursor(c Cursor) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("marshal cursor: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func DecodeCursor(value string) (Cursor, error) {
	if value == "" {
		return Cursor{}, nil
	}
	b, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return Cursor{}, ErrInvalid
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return Cursor{}, ErrInvalid
	}
	if c.ID == "" || c.CreatedAt.IsZero() {
		return Cursor{}, ErrInvalid
	}
	return c, nil
}
func NormalizePageSize(value, defaultValue, maximum int) int {
	if value < 1 {
		return defaultValue
	}
	if value > maximum {
		return maximum
	}
	return value
}
