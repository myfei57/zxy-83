package pos

import (
	"errors"
	"path/filepath"
	"testing"

	"trainwash/internal/store"
)

func TestPersistAndClear(t *testing.T) {
	st := store.New(filepath.Join(t.TempDir(), "pos"))
	tracker := NewTracker(st)
	position := Position{TrainID: "T-1", FrontMM: 5000, LengthMM: 24000, ZeroMM: 100}
	if err := tracker.Persist(position); err != nil {
		t.Fatalf("persist: %v", err)
	}
	if !tracker.Persisted() {
		t.Fatal("position should be persisted")
	}
	if got := tracker.Current(); got != position {
		t.Fatalf("current mismatch: %+v", got)
	}
	if err := tracker.ClearPosition(); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if tracker.Persisted() {
		t.Fatal("position should be cleared")
	}
	if err := tracker.Persist(Position{TrainID: "", FrontMM: 1, LengthMM: 2}); !errors.Is(err, ErrBadPosition) {
		t.Fatalf("expected ErrBadPosition, got %v", err)
	}
}

func TestPointFor(t *testing.T) {
	st := store.New(filepath.Join(t.TempDir(), "pos"))
	tracker := NewTracker(st)
	if err := tracker.Persist(Position{TrainID: "T-2", FrontMM: 3000, LengthMM: 20000, ZeroMM: 0}); err != nil {
		t.Fatalf("persist: %v", err)
	}
	if got := tracker.PointFor(500); got != 3500 {
		t.Fatalf("expected point 3500, got %d", got)
	}
}
