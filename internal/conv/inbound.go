package conv

import (
	"errors"

	"trainwash/internal/store"
)

var ErrAlreadyRunning = errors.New("conv: conveyor already running")

func (s *Service) Inbound() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return ErrAlreadyRunning
	}
	s.running = true
	return store.SaveJSON(s.store, store.KeyConvState, ConvState{Running: true, PositionMM: s.positionMM})
}

func (s *Service) Outbound() error {
	distance := s.layout.NearestWashOffset(s.positionMM) + 8000
	return s.Move(500, distance)
}
