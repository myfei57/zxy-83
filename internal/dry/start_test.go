package dry

import (
	"path/filepath"
	"testing"

	"trainwash/internal/ns"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

func newWaterSystem(t *testing.T) *water.System {
	t.Helper()
	st := store.New(filepath.Join(t.TempDir(), "state"))
	return water.NewSystem(st, ns.DefaultLimits())
}

// The dryer must not start until rinsing has completed. Before the fix the
// guard was inverted: it rejected only StateIdle, so the fans spun up while
// water was still flowing or rinsing was in progress.
func TestStartRejectedBeforeRinseComplete(t *testing.T) {
	w := newWaterSystem(t)
	d := NewService(w)

	if err := w.Start(); err != nil {
		t.Fatalf("water start: %v", err)
	}

	// Flowing, but rinse has not even begun.
	if err := d.Start(); !IsRinseNotDone(err) {
		t.Fatalf("dryer must not start while only flowing, got %v", err)
	}

	// Rinse started but not yet done.
	if err := w.BeginRinse(); err != nil {
		t.Fatalf("begin rinse: %v", err)
	}
	if err := d.Start(); !IsRinseNotDone(err) {
		t.Fatalf("dryer must not start while rinsing, got %v", err)
	}

	// Rinse fully done -> dryer may start.
	if err := w.RinseDone(); err != nil {
		t.Fatalf("rinse done: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("dryer should start after rinse complete, got %v", err)
	}
}

// After an emergency stop the water system is left in StateStopped, which is
// not a rinsed-complete state. Re-running the cycle must not let the fans spin
// up again until a fresh rinse has completed — this is the "复洗一遍后还是提前
// 启动" regression.
func TestStartRejectedAfterStopBeforeRinseRedone(t *testing.T) {
	w := newWaterSystem(t)
	d := NewService(w)

	// Drive a full rinse to completion, then stop the water system.
	_ = w.Start()
	_ = w.BeginRinse()
	if err := w.RinseDone(); err != nil {
		t.Fatalf("rinse done: %v", err)
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if w.RinseComplete() {
		t.Fatal("RinseComplete must be false after stop")
	}
	if err := d.Start(); !IsRinseNotDone(err) {
		t.Fatalf("dryer must not start after stop without re-rinsing, got %v", err)
	}
}

// ResetCycle must clear a leftover running flag. Without it a dryer that was
// running when a cycle was interrupted keeps reporting fans-on into the next
// cycle, so the fans spin before that cycle's rinse completes — the
// "日志里冲洗完成事件出现之前风干机就已经转了" defect.
func TestResetCycleClearsRunningFlag(t *testing.T) {
	w := newWaterSystem(t)
	d := NewService(w)

	_ = w.Start()
	_ = w.BeginRinse()
	if err := w.RinseDone(); err != nil {
		t.Fatalf("rinse done: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("dry start: %v", err)
	}
	if !d.Running() {
		t.Fatal("dryer should be running after Start")
	}

	// A new cycle begins: the cycle boundary stops the water system (as
	// EmergencyStop does), so rinse is no longer complete, and ResetCycle
	// clears the dryer's stale running flag.
	_ = w.Stop()
	d.ResetCycle()
	if d.Running() {
		t.Fatal("ResetCycle must clear the running flag")
	}
	if d.FanSpeed() != 0 {
		t.Fatalf("fans must be off after ResetCycle, got speed=%d", d.FanSpeed())
	}
	// Rinse is not complete for the new cycle, so the guard must reject Start.
	if err := d.Start(); !IsRinseNotDone(err) {
		t.Fatalf("Start must be rejected until a fresh rinse completes, got %v", err)
	}
}

// IsRinseNotDone reports whether err is ErrRinseNotDone.
func IsRinseNotDone(err error) bool { return err == ErrRinseNotDone }
