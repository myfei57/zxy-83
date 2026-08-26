package verifycase

import (
	"path/filepath"
	"testing"

	"trainwash/internal/dry"
	"trainwash/internal/ns"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

func TestTwsDryAfterRinse(t *testing.T) {
	root := filepath.Join(t.TempDir(), "state")
	st := store.New(root)
	limits := ns.DefaultLimits()
	waterSystem := water.NewSystem(st, limits)
	dryService := dry.NewService(waterSystem)
	if err := waterSystem.Start(); err != nil {
		t.Fatalf("water start: %v", err)
	}
	dryService.ResetCycle()
	if err := dryService.Start(); err == nil {
		t.Fatal("dryer must not start before the rinse completes")
	}
	if err := waterSystem.BeginRinse(); err != nil {
		t.Fatalf("begin rinse: %v", err)
	}
	if err := waterSystem.RinseDone(); err != nil {
		t.Fatalf("rinse done: %v", err)
	}
	if err := dryService.Start(); err != nil {
		t.Fatalf("dryer start after rinse: %v", err)
	}
	if err := dryService.Stop(); err != nil {
		t.Fatalf("dryer stop: %v", err)
	}
	restarted := water.NewSystem(st, limits)
	restartedDry := dry.NewService(restarted)
	if err := restartedDry.Start(); err != nil {
		t.Fatalf("restart must restore the rinse state: %v", err)
	}
	if err := restartedDry.Stop(); err != nil {
		t.Fatalf("restarted dryer stop: %v", err)
	}
	if err := waterSystem.Stop(); err != nil {
		t.Fatalf("water stop: %v", err)
	}
	if err := waterSystem.Drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}
	if err := waterSystem.Start(); err != nil {
		t.Fatalf("second water start: %v", err)
	}
	dryService.ResetCycle()
	if err := dryService.Start(); err == nil {
		t.Fatal("re-wash dryer start must wait for the rinse again")
	}
}
