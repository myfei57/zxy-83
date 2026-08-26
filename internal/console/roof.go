package console

import "net/http"

func (s *System) handleRoofLower(w http.ResponseWriter, r *http.Request) {
	if err := s.Roof.Lower(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"lowered": true})
}

func (s *System) handleRoofRaise(w http.ResponseWriter, r *http.Request) {
	if err := s.Roof.Raise(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"lowered": false})
}

func (s *System) handleRoofBrush(w http.ResponseWriter, r *http.Request) {
	swept, err := s.Roof.Brush(0, 1000)
	if err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	point, err := s.Roof.BrushZone(500)
	if err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"swept_mm": swept, "point_mm": point})
}
