package conv

import (
	"errors"

	"trainwash/internal/store"
)

var (
	ErrNotRunning  = errors.New("conv: conveyor is not running")
	ErrBadSpeed    = errors.New("conv: invalid speed")
	ErrBadDistance = errors.New("conv: invalid distance")
)

func (s *Service) Move(speedMMS, distanceMM int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return ErrNotRunning
	}
	if !s.limits.SpeedOK(speedMMS) {
		return ErrBadSpeed
	}
	if distanceMM < 0 {
		return ErrBadDistance
	}
	s.positionMM += distanceMM
	return store.SaveJSON(s.store, store.KeyConvState, ConvState{Running: true, PositionMM: s.positionMM})
}
