package console

import (
	"net/http"

	"trainwash/internal/pos"
)

type positionRequest struct {
	TrainID  string `json:"train_id"`
	FrontMM  int    `json:"front_mm"`
	LengthMM int    `json:"length_mm"`
	ZeroMM   int    `json:"zero_mm"`
}

func (s *System) handlePosPersist(w http.ResponseWriter, r *http.Request) {
	var req positionRequest
	if !decodeBody(w, r, &req) {
		return
	}
	if err := s.Pos.Persist(positionFromRequest(req)); err != nil {
		writeErr(w, http.StatusBadRequest, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"persisted": s.Pos.Persisted(), "position": s.Pos.Current()})
}

func (s *System) handlePosCalibrate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ZeroMM int `json:"zero_mm"`
	}
	if !decodeBody(w, r, &req) {
		return
	}
	if err := s.Pos.Recalibrate(req.ZeroMM); err != nil {
		writeErr(w, http.StatusBadRequest, errMessage(err))
		return
	}
	s.Brush.RefreshZero()
	writeOK(w, map[string]any{"zero_mm": s.Pos.ZeroMM()})
}

func (s *System) handlePosHead(w http.ResponseWriter, r *http.Request) {
	var req struct {
		HeadMM int `json:"head_mm"`
	}
	if !decodeBody(w, r, &req) {
		return
	}
	if err := s.Pos.MarkHead(req.HeadMM); err != nil {
		writeErr(w, http.StatusBadRequest, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"head_arrived": s.Pos.HeadArrived(), "head_mm": s.Pos.HeadMM()})
}

func positionFromRequest(req positionRequest) pos.Position {
	return pos.Position{TrainID: req.TrainID, FrontMM: req.FrontMM, LengthMM: req.LengthMM, ZeroMM: req.ZeroMM}
}
