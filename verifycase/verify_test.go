package verifycase

import (
	"path/filepath"
	"testing"

	"trainwash/internal/audit"
	"trainwash/internal/brush"
	"trainwash/internal/chem"
	"trainwash/internal/conv"
	"trainwash/internal/dry"
	"trainwash/internal/entry"
	"trainwash/internal/ns"
	"trainwash/internal/plan"
	"trainwash/internal/pos"
	"trainwash/internal/roof"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

func buildSystem(t *testing.T) (*store.FileStore, *entry.Service, *brush.Set, *plan.Cycle) {
	t.Helper()
	root := filepath.Join(t.TempDir(), "state")
	st := store.New(root)
	limits := ns.DefaultLimits()
	layout := ns.NewStationLayout()
	tracker := pos.NewTracker(st)
	entryService := entry.NewService(st)
	waterSystem := water.NewSystem(st, limits)
	brushSet := brush.NewSet(st, tracker, waterSystem)
	entryService.AttachPublisher(brushSet)
	chemService := chem.NewService(st, waterSystem, limits)
	roofService := roof.NewService(tracker)
	dryService := dry.NewService(waterSystem)
	convService := conv.NewService(st, waterSystem, layout, limits)
	recorder := audit.NewRecorder(st)
	cycle := plan.NewCycle(entryService, tracker, brushSet, roofService, waterSystem, chemService, dryService, convService, recorder)
	brushSet.AttachReleaser(entryService)
	waterSystem.AttachDrainGate(chemService)
	_ = cycle.Recover()
	return st, entryService, brushSet, cycle
}

func TestTwsBrushLatchReset(t *testing.T) {
	st, entryService, _, cycle := buildSystem(t)
	train1 := entry.NewTrain("T-1", entry.TypeShort, 24000)
	position1 := pos.Position{TrainID: "T-1", FrontMM: 5000, LengthMM: 24000}
	if err := cycle.StartWash(train1, position1); err != nil {
		t.Fatalf("first wash start: %v", err)
	}
	restartedEntry := entry.NewService(st)
	if !restartedEntry.Latched() {
		t.Fatal("restart must keep the entry gate latched during a wash")
	}
	if err := cycle.CompleteWash(); err != nil {
		t.Fatalf("complete wash: %v", err)
	}
	if entryService.Latched() {
		t.Fatal("entry gate latch must reset after the brush arm retracts")
	}
	train2 := entry.NewTrain("T-2", entry.TypeLong, 40000)
	position2 := pos.Position{TrainID: "T-2", FrontMM: 7000, LengthMM: 40000}
	if err := cycle.StartWash(train2, position2); err != nil {
		t.Fatalf("next train must be accepted after retract: %v", err)
	}
}
