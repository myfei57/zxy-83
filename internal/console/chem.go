package console

import "net/http"

func (s *System) handleChemSpray(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DurationMS int `json:"duration_ms"`
	}
	if !decodeBody(w, r, &req) {
		return
	}
	if err := s.Chem.Spray(req.DurationMS); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"dose_ml": s.Chem.LastDoseML()})
}

func (s *System) handleChemAlarm(w http.ResponseWriter, r *http.Request) {
	if err := s.Chem.SetAlarm(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"alarm": s.Chem.AlarmActive(), "valve_latched": s.Chem.ValveLatched()})
}

func (s *System) handleChemAlarmClear(w http.ResponseWriter, r *http.Request) {
	if err := s.Chem.AlarmClear(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"alarm": s.Chem.AlarmActive(), "valve_latched": s.Chem.ValveLatched(), "dose_ml": s.Chem.DoseML()})
}
