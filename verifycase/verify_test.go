package verifycase

import (
	"path/filepath"
	"testing"

	"trainwash/internal/conv"
	"trainwash/internal/ns"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

func TestTwsStopOrder(t *testing.T) {
	root := filepath.Join(t.TempDir(), "state")
	st := store.New(root)
	limits := ns.DefaultLimits()
	layout := ns.NewStationLayout()
	waterSystem := water.NewSystem(st, limits)
	convService := conv.NewService(st, waterSystem, layout, limits)
	if err := waterSystem.Start(); err != nil {
		t.Fatalf("water start: %v", err)
	}
	if err := convService.Stop(); err == nil {
		t.Fatal("conveyor must not stop while the water system is still running")
	}
	if err := waterSystem.Stop(); err != nil {
		t.Fatalf("water stop: %v", err)
	}
	if err := convService.Stop(); err != nil {
		t.Fatalf("conveyor stop after water: %v", err)
	}
	if convService.Running() {
		t.Fatal("conveyor must be stopped")
	}
	restarted := conv.NewService(st, waterSystem, layout, limits)
	if restarted.Running() {
		t.Fatal("restart must restore the stopped conveyor state")
	}
}
