package pos

import (
	"path/filepath"
	"testing"

	"trainwash/internal/store"
)

func TestMarkHeadSetsArrived(t *testing.T) {
	st := store.New(filepath.Join(t.TempDir(), "pos"))
	tracker := NewTracker(st)
	if tracker.HeadArrived() {
		t.Fatal("head should not be arrived before MarkHead")
	}
	if err := tracker.MarkHead(8000); err != nil {
		t.Fatalf("mark head: %v", err)
	}
	if !tracker.HeadArrived() {
		t.Fatal("head should be arrived after MarkHead")
	}
	if got := tracker.HeadMM(); got != 8000 {
		t.Fatalf("expected head mm 8000, got %d", got)
	}
}

func TestClearHeadClearsArrived(t *testing.T) {
	st := store.New(filepath.Join(t.TempDir(), "pos"))
	tracker := NewTracker(st)
	if err := tracker.MarkHead(8000); err != nil {
		t.Fatalf("mark head: %v", err)
	}
	if err := tracker.ClearHead(); err != nil {
		t.Fatalf("clear head: %v", err)
	}
	if tracker.HeadArrived() {
		t.Fatal("head should not be arrived after ClearHead")
	}
}

// TestMarkHeadRestoredAcrossRestart locks in the regression: before the fix,
// restoreHead was empty, so HeadArrived() always read false after a restart
// even when a head had been persisted to disk.
func TestMarkHeadRestoredAcrossRestart(t *testing.T) {
	dir := t.TempDir()
	st := store.New(filepath.Join(dir, "pos"))
	tracker := NewTracker(st)
	if err := tracker.MarkHead(8200); err != nil {
		t.Fatalf("mark head: %v", err)
	}

	// A fresh tracker over the same store must see the head as arrived.
	reopened := NewTracker(store.New(filepath.Join(dir, "pos")))
	if !reopened.HeadArrived() {
		t.Fatal("head should be arrived after restart")
	}
	if got := reopened.HeadMM(); got != 8200 {
		t.Fatalf("expected head mm 8200 after restart, got %d", got)
	}
}
