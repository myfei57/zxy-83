package water

import (
	"path/filepath"
	"testing"

	"trainwash/internal/ns"
	"trainwash/internal/store"
)

func newSystem(t *testing.T) *System {
	t.Helper()
	st := store.New(filepath.Join(t.TempDir(), "state"))
	return NewSystem(st, ns.DefaultLimits())
}

// RinseDone must drive the system into the StateRinseDone state so that the
// dryer can gate on it. Before the fix it persisted StateRinsing unchanged,
// leaving RinseComplete() false forever.
func TestRinseDoneTransitionsToRinseDone(t *testing.T) {
	s := newSystem(t)
	if err := s.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := s.BeginRinse(); err != nil {
		t.Fatalf("begin rinse: %v", err)
	}
	if err := s.RinseDone(); err != nil {
		t.Fatalf("rinse done: %v", err)
	}
	if got := s.State(); got != StateRinseDone {
		t.Fatalf("state after RinseDone = %s, want %s", got, StateRinseDone)
	}
	if !s.RinseComplete() {
		t.Fatal("RinseComplete() must be true after RinseDone")
	}
}

// Re-running RinseDone after completion is a no-op (already past rinsing).
func TestRinseDoneIdempotentAfterComplete(t *testing.T) {
	s := newSystem(t)
	_ = s.Start()
	_ = s.BeginRinse()
	if err := s.RinseDone(); err != nil {
		t.Fatalf("first rinse done: %v", err)
	}
	if err := s.RinseDone(); !IsErrNotRinsing(err) {
		t.Fatalf("second RinseDone should be a no-op error, got %v", err)
	}
	if got := s.State(); got != StateRinseDone {
		t.Fatalf("state should remain %s, got %s", StateRinseDone, got)
	}
}

func IsErrNotRinsing(err error) bool { return err == ErrNotRinsing }
