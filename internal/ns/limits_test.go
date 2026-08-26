package ns

import "testing"

func TestLimitsPressureOK(t *testing.T) {
	limits := DefaultLimits()
	if !limits.PressureOK(5.0) {
		t.Fatal("5.0 should be within pressure range")
	}
	if limits.PressureOK(20.0) {
		t.Fatal("20.0 should be rejected")
	}
	if limits.PressureOK(0.5) {
		t.Fatal("0.5 should be rejected")
	}
}

func TestLimitsDoseAndSpeed(t *testing.T) {
	limits := DefaultLimits()
	if !limits.DoseOK(30.0) {
		t.Fatal("30ml dose should be accepted")
	}
	if limits.DoseOK(200.0) {
		t.Fatal("200ml dose should be rejected")
	}
	if !limits.SpeedOK(600) {
		t.Fatal("600mm/s should be accepted")
	}
	if limits.SpeedOK(5000) {
		t.Fatal("5000mm/s should be rejected")
	}
	if !limits.SprayDurationOK(15000) {
		t.Fatal("15000ms spray should be accepted")
	}
	if limits.SprayDurationOK(-1) {
		t.Fatal("negative spray should be rejected")
	}
}

func TestSegmentsOf(t *testing.T) {
	segments := SegmentsOf(24000)
	if len(segments) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segments))
	}
	if segments[0].Kind != SegmentHead || segments[2].Kind != SegmentTail {
		t.Fatalf("unexpected segment order: %v", segments)
	}
	if segments[0].Contains(10) == segments[0].Contains(segments[0].EndMM) {
		t.Fatal("segment contains boundary mismatch")
	}
}

func TestStationLayout(t *testing.T) {
	layout := NewStationLayout()
	if layout.Count() != 5 {
		t.Fatalf("expected 5 stations, got %d", layout.Count())
	}
	station, ok := layout.StationByKind(StationWash)
	if !ok || station.OffsetMM != 8000 {
		t.Fatalf("unexpected wash station: %+v ok=%v", station, ok)
	}
	if layout.NearestWashOffset(8200) != -200 {
		t.Fatalf("unexpected wash offset: %d", layout.NearestWashOffset(8200))
	}
}
