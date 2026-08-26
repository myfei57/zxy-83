package verifycase

import (
	"path/filepath"
	"testing"

	"trainwash/internal/brush"
	"trainwash/internal/entry"
	"trainwash/internal/ns"
	"trainwash/internal/pos"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

func TestTwsWaterBaselineFresh(t *testing.T) {
	root := filepath.Join(t.TempDir(), "state")
	st := store.New(root)
	limits := ns.DefaultLimits()
	tracker := pos.NewTracker(st)
	waterSystem := water.NewSystem(st, limits)
	brushSet := brush.NewSet(st, tracker, waterSystem)
	entryService := entry.NewService(st)
	entryService.AttachPublisher(brushSet)
	if err := tracker.Persist(pos.Position{TrainID: "T-1", FrontMM: 5000, LengthMM: 24000}); err != nil {
		t.Fatalf("persist: %v", err)
	}
	if err := brushSet.LowerSide(); err != nil {
		t.Fatalf("lower: %v", err)
	}
	if err := waterSystem.RecalibratePressure(9.0); err != nil {
		t.Fatalf("recalibrate: %v", err)
	}
	mpa, err := brushSet.Rinse(500)
	if err != nil {
		t.Fatalf("rinse: %v", err)
	}
	if mpa != 4.5 {
		t.Fatalf("rinse must use the recalibrated gain, got %v", mpa)
	}
	restarted := water.NewSystem(st, limits)
	if got := restarted.GainMPA(); got != 9.0 {
		t.Fatalf("gain must survive restart, got %v", got)
	}
}
