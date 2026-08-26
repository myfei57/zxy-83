package console

import (
	"net/http"
	"time"

	"trainwash/internal/brush"
	"trainwash/internal/ns"
	"trainwash/internal/store"
)

func (s *System) handleState(w http.ResponseWriter, r *http.Request) {
	position := s.Pos.Current()
	headRecord, _ := s.Pos.HeadRecord()
	station, stationOK := s.Layout.StationAt(position.FrontMM)
	stationInfo := map[string]any{"ok": stationOK}
	if stationOK {
		stationInfo["id"] = station.ID
		stationInfo["kind"] = station.Kind.String()
		stationInfo["offset_mm"] = station.OffsetMM
	}
	writeOK(w, map[string]any{
		"position":       position,
		"zero_mm":        s.Pos.ZeroMM(),
		"persisted":      s.Pos.Persisted(),
		"head_arrived":   s.Pos.HeadArrived(),
		"head_epoch":     s.Pos.HeadEpoch(),
		"head_record":    headRecord,
		"segments":       s.segmentSnapshot(position.LengthMM),
		"station":        stationInfo,
		"data_dir":       s.DataDir,
		"store_position": s.Store.Exists(store.KeyPosPosition),
		"water": map[string]any{
			"state":    s.Water.State().String(),
			"gain_mpa": s.Water.GainMPA(),
			"cycle_id": s.Water.CycleID(),
			"rinse_ok": s.Water.RinseComplete(),
			"stop_ok":  s.Water.StopReady(),
		},
		"brush": map[string]any{
			"active_group":      s.Brush.ActiveGroup(),
			"group_zones":       groupSnapshot(s.Brush.ActiveGroup(), position.FrontMM),
			"lowered":           s.Brush.IsLowered(),
			"zero_mm":           s.Brush.CachedZeroMM(),
			"last_pressure_mpa": s.Brush.LastPressureMPA(),
		},
		"roof": map[string]any{
			"lowered":    s.Roof.IsLowered(),
			"head_ready": s.Roof.HeadReady(),
			"head_mm":    s.Roof.HeadPositionMM(),
		},
		"dry": map[string]any{
			"running":   s.Dry.Running(),
			"fan_speed": s.Dry.FanSpeed(),
			"duty":      s.Dry.DutyCycle(),
		},
		"chem": map[string]any{
			"valve_latched": s.Chem.ValveLatched(),
			"alarm":         s.Chem.AlarmActive(),
			"alarm_seq":     s.Chem.AlarmSeq(),
			"dose_ml":       s.Chem.DoseML(),
			"dose_ok":       s.Chem.DoseOK(),
		},
		"conv": map[string]any{
			"running":     s.Conv.Running(),
			"position_mm": s.Conv.PositionMM(),
		},
		"entry": map[string]any{
			"train_type": s.Entry.CurrentType(),
			"latched":    s.Entry.Latched(),
			"wash_seq":   s.Entry.WashSeq(),
		},
		"cycle": map[string]any{
			"stage":    s.Cycle.Stage().String(),
			"train_id": s.Cycle.TrainID(),
		},
		"layout": map[string]any{
			"stations": s.Layout.Count(),
			"string":   s.Layout.String(),
		},
		"limits": s.Limits,
		"time":   s.Now().UTC().Format(time.RFC3339),
	})
}

func (s *System) segmentSnapshot(lengthMM int) []map[string]any {
	position := s.Pos.Current()
	segments := ns.SegmentsOf(lengthMM)
	out := make([]map[string]any, 0, len(segments))
	for _, seg := range segments {
		out = append(out, map[string]any{
			"kind":           seg.Kind.String(),
			"start_mm":       seg.StartMM,
			"end_mm":         seg.EndMM,
			"length_mm":      seg.LengthMM(),
			"contains_front": seg.Contains(position.FrontMM),
		})
	}
	return out
}

func groupSnapshot(group brush.Group, frontMM int) []map[string]any {
	out := make([]map[string]any, 0, len(group.Zones))
	for _, zone := range group.Zones {
		out = append(out, map[string]any{
			"kind":           zone.Kind.String(),
			"start_mm":       zone.StartMM,
			"end_mm":         zone.EndMM,
			"contains_front": zone.Contains(frontMM),
		})
	}
	return out
}

func (s *System) handleAudit(w http.ResponseWriter, r *http.Request) {
	limit := 100
	records, err := s.Audit.List(limit)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, errMessage(err))
		return
	}
	count, _ := s.Audit.Count()
	writeOK(w, map[string]any{"count": count, "records": records})
}

func (s *System) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func (s *System) handlePage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/"+name)
	}
}
