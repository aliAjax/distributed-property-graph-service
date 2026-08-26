package importer

import (
	"context"
	"errors"
	"testing"

	"github.com/example/distributed-property-graph/internal/platform"
)

func TestSaveRejectsInvalidImportID(t *testing.T) {
	store := NewCheckpointStore()
	err := store.Save(context.Background(), Checkpoint{ImportID: "", Line: 1})
	if !errors.Is(err, platform.ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestGetMissingReturnsNotFound(t *testing.T) {
	store := NewCheckpointStore()
	_, err := store.Get(context.Background(), "missing")
	if !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteMissingReturnsNotFound(t *testing.T) {
	store := NewCheckpointStore()
	err := store.Delete(context.Background(), "missing")
	if !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteExistingRemovesCheckpoint(t *testing.T) {
	store := NewCheckpointStore()
	if err := store.Save(context.Background(), Checkpoint{ImportID: "i", Line: 1}); err != nil {
		t.Fatal(err)
	}
	if err := store.Delete(context.Background(), "i"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(context.Background(), "i"); !errors.Is(err, platform.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}
