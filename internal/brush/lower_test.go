package brush

import (
	"errors"
	"path/filepath"
	"testing"

	"trainwash/internal/ns"
	"trainwash/internal/pos"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

func newBrushSet(t *testing.T) (*Set, *pos.Tracker) {
	t.Helper()
	st := store.New(filepath.Join(t.TempDir(), "brush"))
	tracker := pos.NewTracker(st)
	waterSystem := water.NewSystem(st, ns.DefaultLimits())
	return NewSet(st, tracker, waterSystem), tracker
}

func mustPersistPosition(t *testing.T, tracker *pos.Tracker) {
	t.Helper()
	if err := tracker.Persist(pos.Position{TrainID: "T-1", FrontMM: 5000, LengthMM: 24000, ZeroMM: 0}); err != nil {
		t.Fatalf("persist position: %v", err)
	}
}

// TestLowerSideRequiresHeadArrived locks in the reported incident: the side
// brush dropped onto the windshield because LowerSide() ignored the head
// position. Now lowering before the head has arrived must be refused.
func TestLowerSideRequiresHeadArrived(t *testing.T) {
	set, tracker := newBrushSet(t)
	mustPersistPosition(t, tracker)
	if err := set.PublishGroup("standard", 24000, 600); err != nil {
		t.Fatalf("publish group: %v", err)
	}
	if err := set.LowerSide(); !errors.Is(err, ErrHeadNotReady) {
		t.Fatalf("expected ErrHeadNotReady before head arrives, got %v", err)
	}
	if set.IsLowered() {
		t.Fatal("brush must not be lowered when head has not arrived")
	}
}

func TestLowerSideAllowedAfterHeadArrived(t *testing.T) {
	set, tracker := newBrushSet(t)
	mustPersistPosition(t, tracker)
	if err := tracker.MarkHead(8000); err != nil {
		t.Fatalf("mark head: %v", err)
	}
	if err := set.PublishGroup("standard", 24000, 600); err != nil {
		t.Fatalf("publish group: %v", err)
	}
	if err := set.LowerSide(); err != nil {
		t.Fatalf("lower after head arrived: %v", err)
	}
	if !set.IsLowered() {
		t.Fatal("brush should be lowered after LowerSide succeeds")
	}
	if err := set.LowerSide(); !errors.Is(err, ErrAlreadyLowered) {
		t.Fatalf("expected ErrAlreadyLowered on double lower, got %v", err)
	}
}

func TestLowerSideStillRequiresPersistedAndGroup(t *testing.T) {
	set, tracker := newBrushSet(t)
	// No position persisted, no head, no group.
	if err := set.LowerSide(); !errors.Is(err, ErrNotPersisted) {
		t.Fatalf("expected ErrNotPersisted, got %v", err)
	}
	// Persist + head, but still no active group.
	mustPersistPosition(t, tracker)
	if err := tracker.MarkHead(8000); err != nil {
		t.Fatalf("mark head: %v", err)
	}
	if err := set.LowerSide(); !errors.Is(err, ErrNoGroup) {
		t.Fatalf("expected ErrNoGroup, got %v", err)
	}
}
