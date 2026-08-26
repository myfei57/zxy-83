package console

import "net/http"

func (s *System) handleBrushLower(w http.ResponseWriter, r *http.Request) {
	if err := s.Brush.LowerSide(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"lowered": true, "zero_mm": s.Brush.CachedZeroMM()})
}

func (s *System) handleBrushRetract(w http.ResponseWriter, r *http.Request) {
	if err := s.Brush.Retract(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"lowered": false})
}

func (s *System) handleBrushRaise(w http.ResponseWriter, r *http.Request) {
	if err := s.Brush.RaiseSide(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"lowered": false})
}
