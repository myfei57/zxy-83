package entry

import (
	"path/filepath"
	"testing"

	"trainwash/internal/store"
)

func TestLatchReleaseClearsGate(t *testing.T) {
	s := NewService(store.New(filepath.Join(t.TempDir(), "entry")))

	if s.Latched() {
		t.Fatal("latch should start cleared")
	}
	if err := s.Latch(); err != nil {
		t.Fatalf("latch: %v", err)
	}
	if !s.Latched() {
		t.Fatal("latch should be set after Latch")
	}

	// A second train must be rejected while the gate is still latched.
	if err := s.Latch(); !errIsLatched(err) {
		t.Fatalf("expected ErrLatched, got %v", err)
	}

	// Releasing must clear the in-memory flag so the next train can enter.
	if err := s.ReleaseLatch(); err != nil {
		t.Fatalf("release: %v", err)
	}
	if s.Latched() {
		t.Fatal("latch should be cleared after ReleaseLatch")
	}
	if err := s.Latch(); err != nil {
		t.Fatalf("latch after release: %v", err)
	}
}

func TestLatchRestoredFromStore(t *testing.T) {
	root := filepath.Join(t.TempDir(), "entry")
	st := store.New(root)

	first := NewService(st)
	if err := first.Latch(); err != nil {
		t.Fatalf("latch: %v", err)
	}
	seq := first.WashSeq()

	// Reconstruct from the same store as if the process restarted: the
	// latched state must survive so a crash mid-wash does not admit the
	// next train into an occupied bay.
	second := NewService(st)
	if !second.Latched() {
		t.Fatal("latch should be restored from store")
	}
	if second.WashSeq() != seq {
		t.Fatalf("wash seq mismatch: got %d want %d", second.WashSeq(), seq)
	}
}

func errIsLatched(err error) bool {
	return err != nil && err.Error() == ErrLatched.Error()
}
