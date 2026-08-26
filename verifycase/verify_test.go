package verifycase

import (
	"path/filepath"
	"testing"

	"trainwash/internal/pos"
	"trainwash/internal/roof"
	"trainwash/internal/store"
)

func TestTwsRoofAfterHead(t *testing.T) {
	root := filepath.Join(t.TempDir(), "state")
	st := store.New(root)
	tracker := pos.NewTracker(st)
	roofService := roof.NewService(tracker)
	if err := roofService.Lower(); err == nil {
		t.Fatal("roof brush must not lower before the train head arrives")
	}
	if err := tracker.MarkHead(3000); err != nil {
		t.Fatalf("mark head: %v", err)
	}
	if err := roofService.Lower(); err != nil {
		t.Fatalf("lower after head: %v", err)
	}
	if err := roofService.Raise(); err != nil {
		t.Fatalf("raise: %v", err)
	}
	restarted := pos.NewTracker(st)
	restartedRoof := roof.NewService(restarted)
	if err := restartedRoof.Lower(); err != nil {
		t.Fatalf("restart must restore the head record: %v", err)
	}
	if err := restartedRoof.Raise(); err != nil {
		t.Fatalf("raise after restart: %v", err)
	}
	if err := restarted.ClearHead(); err != nil {
		t.Fatalf("clear head: %v", err)
	}
	if err := restartedRoof.Lower(); err == nil {
		t.Fatal("re-entry roof lower must wait for the head again")
	}
	if err := restarted.MarkHead(3500); err != nil {
		t.Fatalf("mark head again: %v", err)
	}
	if err := restartedRoof.Lower(); err != nil {
		t.Fatalf("lower after second head: %v", err)
	}
}
