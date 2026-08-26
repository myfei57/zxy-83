package verifycase

import (
	"path/filepath"
	"testing"

	"trainwash/internal/chem"
	"trainwash/internal/ns"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

func TestTwsChemLatchRelease(t *testing.T) {
	root := filepath.Join(t.TempDir(), "state")
	st := store.New(root)
	limits := ns.DefaultLimits()
	waterSystem := water.NewSystem(st, limits)
	chemService := chem.NewService(st, waterSystem, limits)
	waterSystem.AttachDrainGate(chemService)
	if err := chemService.SetAlarm(); err != nil {
		t.Fatalf("set alarm: %v", err)
	}
	restarted := chem.NewService(st, waterSystem, limits)
	if !restarted.ValveLatched() {
		t.Fatal("restart must restore the latched proportioning valve")
	}
	if err := chemService.AlarmClear(); err != nil {
		t.Fatalf("alarm clear: %v", err)
	}
	if chemService.AlarmActive() {
		t.Fatal("alarm must clear")
	}
	if chemService.ValveLatched() {
		t.Fatal("valve latch must release when the alarm clears")
	}
	if got := chemService.DoseML(); got != 24.0 {
		t.Fatalf("dose must return to normal, got %v", got)
	}
	if err := waterSystem.Start(); err != nil {
		t.Fatalf("water start: %v", err)
	}
	if err := waterSystem.BeginRinse(); err != nil {
		t.Fatalf("begin rinse: %v", err)
	}
	if err := waterSystem.RinseDone(); err != nil {
		t.Fatalf("rinse done: %v", err)
	}
	if err := waterSystem.Drain(); err != nil {
		t.Fatalf("drain must succeed after alarm recovery: %v", err)
	}
}
