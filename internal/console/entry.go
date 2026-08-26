package console

import (
	"net/http"

	"trainwash/internal/entry"
)

func (s *System) handleEntryType(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type string `json:"type"`
	}
	if !decodeBody(w, r, &req) {
		return
	}
	trainType := entry.NormalizeType(req.Type)
	if err := s.Entry.TypeChange(trainType); err != nil {
		writeErr(w, http.StatusBadRequest, errMessage(err))
		return
	}
	writeOK(w, map[string]any{
		"train_type": s.Entry.CurrentType(),
		"group":      s.Brush.ActiveGroup(),
		"expected":   s.Entry.ResolveGroup(trainType),
	})
}

func (s *System) handleConvInbound(w http.ResponseWriter, r *http.Request) {
	if err := s.Conv.Inbound(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"running": s.Conv.Running()})
}

func (s *System) handleConvMove(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SpeedMMS   int `json:"speed_mms"`
		DistanceMM int `json:"distance_mm"`
	}
	if !decodeBody(w, r, &req) {
		return
	}
	if err := s.Conv.Move(req.SpeedMMS, req.DistanceMM); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"position_mm": s.Conv.PositionMM()})
}

func (s *System) handleConvStop(w http.ResponseWriter, r *http.Request) {
	if err := s.Conv.Stop(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"running": false})
}
