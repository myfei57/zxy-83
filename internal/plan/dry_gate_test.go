package plan

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
	"trainwash/internal/pos"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

func newCycle(t *testing.T) (*Cycle, *water.System, *dry.Service) {
	t.Helper()
	st := store.New(filepath.Join(t.TempDir(), "state"))
	layout := ns.NewStationLayout()
	limits := ns.DefaultLimits()
	tracker := pos.NewTracker(st)
	entrySvc := entry.NewService(st)
	waterSys := water.NewSystem(st, limits)
	brushSet := brush.NewSet(st, tracker, waterSys)
	entrySvc.AttachPublisher(brushSet)
	chemSvc := chem.NewService(st, waterSys, limits)
	drySvc := dry.NewService(waterSys)
	convSvc := conv.NewService(st, waterSys, layout, limits)
	recorder := audit.NewRecorder(st)
	c := NewCycle(entrySvc, tracker, brushSet, nil, waterSys, chemSvc, drySvc, convSvc, recorder)
	return c, waterSys, drySvc
}

// Regression for the field report: "复洗一遍后还是提前启动，日志里冲洗完成事件
// 出现之前风干机就已经转了". After a cycle is interrupted mid-dry (emergency
// stop) the dryer's running flag is never cleared — EmergencyStop stops water
// and conveyor but never the dryer, and ResetCycle is a no-op. On re-wash the
// fans report running (FanSpeed 1800) BEFORE water.RinseDone() fires, which is
// exactly the watermark-on-bodywork defect described.
func TestDryerFanNotSpinningBeforeRinseDoneOnRewash(t *testing.T) {
	c, waterSys, drySvc := newCycle(t)

	train := entry.NewTrain("T-1", entry.TypeShort, 20000)
	posn := pos.Position{TrainID: "T-1", FrontMM: 100, LengthMM: 20000, ZeroMM: 0}
	if err := c.StartWash(train, posn); err != nil {
		t.Fatalf("start wash: %v", err)
	}
	// Drive a normal rinse+dry so the dryer's running flag is set true.
	if err := waterSys.BeginRinse(); err != nil {
		t.Fatalf("begin rinse: %v", err)
	}
	if err := waterSys.RinseDone(); err != nil {
		t.Fatalf("rinse done: %v", err)
	}
	if err := drySvc.Start(); err != nil {
		t.Fatalf("dry start: %v", err)
	}
	if !drySvc.Running() {
		t.Fatal("dryer should be running after Start")
	}

	// Emergency stop interrupts mid-dry: water + conveyor stopped. EmergencyStop
	// does not stop the dryer (it should), so the fan keeps reporting running.
	if err := c.EmergencyStop(); err != nil {
		t.Fatalf("emergency stop: %v", err)
	}
	// Release the entry gate latch and retract brushes so a fresh wash can
	// begin — these are independent cleanup steps, not the dryer regression.
	_ = c.entry.ReleaseLatch()
	_ = c.brush.Retract()

	// Re-wash: a fresh cycle begins. Rinse has NOT completed yet.
	if err := c.StartWash(train, posn); err != nil {
		t.Fatalf("re-wash: %v", err)
	}
	if waterSys.RinseComplete() {
		t.Fatal("rinse must not be complete at start of re-wash")
	}
	// The fan must NOT be turning before RinseDone on the new cycle.
	if drySvc.Running() || drySvc.FanSpeed() != 0 {
		t.Fatalf("REGRESSION: fan spinning before RinseDone on re-wash (running=%v fan_speed=%d water=%s)",
			drySvc.Running(), drySvc.FanSpeed(), waterSys.State())
	}
}
