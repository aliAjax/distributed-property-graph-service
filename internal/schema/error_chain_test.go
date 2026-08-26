package schema

import (
	"context"
	"errors"
	"testing"

	"github.com/example/distributed-property-graph/internal/platform"
)

func TestValidatePropertyInvalidTypeWrapsSentinel(t *testing.T) {
	err := ValidateProperty(Property{Name: "p", Type: "unknown"})
	if !errors.Is(err, platform.ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

type failingRepo struct{}

func (failingRepo) Save(context.Context, Schema) error { return platform.ErrInvalid }
func (failingRepo) Get(context.Context, string) (Schema, error) { return Schema{}, platform.ErrNotFound }

func TestPublishSaveErrorWrapsSentinel(t *testing.T) {
	svc := NewService(failingRepo{})
	value := Schema{GraphID: "g", Vertices: []VertexType{{Name: "person"}}}
	_, err := svc.Publish(context.Background(), value)
	if !errors.Is(err, platform.ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestSchemaGetMissingWrapsSentinel(t *testing.T) {
	_, err := NewMemoryRepository().Get(context.Background(), "missing")
	if !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestValidateSchemaMissingGraphWrapsSentinel(t *testing.T) {
	err := ValidateSchema(Schema{})
	if !errors.Is(err, platform.ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

type conflictRepo struct{}

func (conflictRepo) Save(context.Context, Schema) error { return nil }
func (conflictRepo) Get(context.Context, string) (Schema, error) { return Schema{}, platform.ErrConflict }

func TestPublishGetErrorPropagates(t *testing.T) {
	svc := NewService(conflictRepo{})
	value := Schema{GraphID: "g", Vertices: []VertexType{{Name: "person"}}}
	if _, err := svc.Publish(context.Background(), value); !errors.Is(err, platform.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
