package console

import (
	"net/http"

	"trainwash/internal/entry"
	"trainwash/internal/pos"
)

type washRequest struct {
	TrainID  string `json:"train_id"`
	Type     string `json:"type"`
	LengthMM int    `json:"length_mm"`
	FrontMM  int    `json:"front_mm"`
	ZeroMM   int    `json:"zero_mm"`
}

func (s *System) handlePlanWash(w http.ResponseWriter, r *http.Request) {
	var req washRequest
	if !decodeBody(w, r, &req) {
		return
	}
	train := entry.NewTrain(req.TrainID, entry.NormalizeType(req.Type), req.LengthMM)
	position := pos.Position{TrainID: req.TrainID, FrontMM: req.FrontMM, LengthMM: req.LengthMM, ZeroMM: req.ZeroMM}
	if err := s.Cycle.StartWash(train, position); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"stage": s.Cycle.Stage().String(), "train_id": s.Cycle.TrainID()})
}

func (s *System) handlePlanComplete(w http.ResponseWriter, r *http.Request) {
	if err := s.Cycle.CompleteWash(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"stage": s.Cycle.Stage().String()})
}

func (s *System) handlePlanStop(w http.ResponseWriter, r *http.Request) {
	if err := s.Cycle.EmergencyStop(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"stage": s.Cycle.Stage().String()})
}
