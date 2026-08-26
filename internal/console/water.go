package console

import (
	"net/http"

	"trainwash/internal/water"
)

func (s *System) handleWaterStart(w http.ResponseWriter, r *http.Request) {
	if !s.Water.CanStart() {
		writeErr(w, http.StatusConflict, errMessage(water.ErrAlreadyRunning))
		return
	}
	if err := s.Water.Start(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"state": s.Water.State().String()})
}

func (s *System) handleWaterStop(w http.ResponseWriter, r *http.Request) {
	if err := s.Water.Stop(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"state": s.Water.State().String()})
}

func (s *System) handleWaterRinse(w http.ResponseWriter, r *http.Request) {
	if err := s.Water.BeginRinse(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	if err := s.Water.RinseDone(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"state": s.Water.State().String(), "rinse_ok": s.Water.RinseComplete()})
}

type recalibrateRequest struct {
	GainMPA float64 `json:"gain_mpa"`
	ZeroMM  int     `json:"zero_mm"`
}

func (s *System) handleWaterRecalibrate(w http.ResponseWriter, r *http.Request) {
	var req recalibrateRequest
	if !decodeBody(w, r, &req) {
		return
	}
	if err := s.Water.RecalibratePressure(req.GainMPA); err != nil {
		writeErr(w, http.StatusBadRequest, errMessage(err))
		return
	}
	s.Brush.RefreshGain()
	expected, err := s.Water.ExpectedPressure(s.Water.GainMPA(), 500)
	if err != nil {
		writeErr(w, http.StatusBadRequest, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"gain_mpa": s.Water.GainMPA(), "expected_mpa_500": expected})
}

func (s *System) handleWaterDrain(w http.ResponseWriter, r *http.Request) {
	if err := s.Water.Drain(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"state": s.Water.State().String()})
}
