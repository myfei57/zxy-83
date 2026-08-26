package brush

import (
	"path/filepath"
	"testing"

	"trainwash/internal/entry"
	"trainwash/internal/ns"
	"trainwash/internal/pos"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

// TestRetractReleasesEntryLatch is the regression for the reported outage:
// after the side brushes retract, the entry latch must be released so the
// next train is not rejected at the gate.
func TestRetractReleasesEntryLatch(t *testing.T) {
	root := filepath.Join(t.TempDir(), "wash")
	st := store.New(root)

	tracker := pos.NewTracker(st)
	if err := tracker.Persist(pos.Position{
		TrainID:  "T-1",
		FrontMM:  5000,
		LengthMM: 24000,
		ZeroMM:   100,
	}); err != nil {
		t.Fatalf("persist position: %v", err)
	}

	limits := ns.DefaultLimits()
	waterSystem := water.NewSystem(st, limits)
	brushSet := NewSet(st, tracker, waterSystem)
	entryService := entry.NewService(st)
	// Mirrors wiring.go: the brush set releases the entry latch on retract.
	brushSet.AttachReleaser(entryService)

	if err := brushSet.ApplyGroup(GroupFromSpec("short-set", 20000, 600)); err != nil {
		t.Fatalf("apply group: %v", err)
	}
	if err := brushSet.LowerSide(); err != nil {
		t.Fatalf("lower side: %v", err)
	}
	if err := entryService.Latch(); err != nil {
		t.Fatalf("latch entry: %v", err)
	}
	if !entryService.Latched() {
		t.Fatal("entry should be latched during wash")
	}

	// Retracting the brushes is the action that previously left the latch set.
	if err := brushSet.Retract(); err != nil {
		t.Fatalf("retract: %v", err)
	}
	if brushSet.IsLowered() {
		t.Fatal("brush should be raised after retract")
	}
	if entryService.Latched() {
		t.Fatal("entry latch must be released after retract so next train can enter")
	}

	// The next train must now be admitted.
	if err := entryService.Latch(); err != nil {
		t.Fatalf("latch after retract: %v", err)
	}
}

// TestResetLatchWithoutReleaser keeps Retract usable when no releaser is
// attached (e.g. unit wiring without the entry service): it must be a no-op,
// not a nil-dereference panic.
func TestResetLatchWithoutReleaser(t *testing.T) {
	root := filepath.Join(t.TempDir(), "wash")
	st := store.New(root)

	tracker := pos.NewTracker(st)
	if err := tracker.Persist(pos.Position{
		TrainID:  "T-1",
		FrontMM:  1000,
		LengthMM: 20000,
		ZeroMM:   0,
	}); err != nil {
		t.Fatalf("persist position: %v", err)
	}
	brushSet := NewSet(st, tracker, water.NewSystem(st, ns.DefaultLimits()))
	if err := brushSet.ApplyGroup(GroupFromSpec("short-set", 20000, 600)); err != nil {
		t.Fatalf("apply group: %v", err)
	}
	if err := brushSet.LowerSide(); err != nil {
		t.Fatalf("lower side: %v", err)
	}
	if err := brushSet.Retract(); err != nil {
		t.Fatalf("retract without releaser: %v", err)
	}
}
