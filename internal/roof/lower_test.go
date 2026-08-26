package roof

import (
	"errors"
	"path/filepath"
	"testing"

	"trainwash/internal/pos"
	"trainwash/internal/store"
)

func newRoofTracker(t *testing.T) (*Service, *pos.Tracker) {
	t.Helper()
	st := store.New(filepath.Join(t.TempDir(), "roof"))
	tracker := pos.NewTracker(st)
	return NewService(tracker), tracker
}

// TestLowerRequiresHeadArrived locks in the reported incident: the roof brush
// dropped onto the windshield because Lower() ignored the head position. Now
// lowering before the head has arrived must be refused.
func TestLowerRequiresHeadArrived(t *testing.T) {
	roof, tracker := newRoofTracker(t)
	_ = tracker
	if err := roof.Lower(); !errors.Is(err, ErrHeadNotReady) {
		t.Fatalf("expected ErrHeadNotReady before head arrives, got %v", err)
	}
	if roof.IsLowered() {
		t.Fatal("roof must not be lowered when head has not arrived")
	}
}

func TestLowerAllowedAfterHeadArrived(t *testing.T) {
	roof, tracker := newRoofTracker(t)
	if err := tracker.MarkHead(8000); err != nil {
		t.Fatalf("mark head: %v", err)
	}
	if err := roof.Lower(); err != nil {
		t.Fatalf("lower after head arrived: %v", err)
	}
	if !roof.IsLowered() {
		t.Fatal("roof should be lowered after Lower succeeds")
	}
	if err := roof.Lower(); !errors.Is(err, ErrAlreadyLowered) {
		t.Fatalf("expected ErrAlreadyLowered on double lower, got %v", err)
	}
}

func TestRaiseRoundTrip(t *testing.T) {
	roof, tracker := newRoofTracker(t)
	if err := tracker.MarkHead(0); err != nil {
		t.Fatalf("mark head: %v", err)
	}
	if err := roof.Lower(); err != nil {
		t.Fatalf("lower: %v", err)
	}
	if err := roof.Raise(); err != nil {
		t.Fatalf("raise: %v", err)
	}
	if roof.IsLowered() {
		t.Fatal("roof should be raised after Raise")
	}
	if err := roof.Raise(); !errors.Is(err, ErrNotLowered) {
		t.Fatalf("expected ErrNotLowered on double raise, got %v", err)
	}
}
