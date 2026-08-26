package console

import "net/http"

func (s *System) handleDryStart(w http.ResponseWriter, r *http.Request) {
	if err := s.Dry.Start(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"running": true, "fan_speed": s.Dry.FanSpeed()})
}

func (s *System) handleDryStop(w http.ResponseWriter, r *http.Request) {
	if err := s.Dry.Stop(); err != nil {
		writeErr(w, http.StatusConflict, errMessage(err))
		return
	}
	writeOK(w, map[string]any{"running": false})
}
