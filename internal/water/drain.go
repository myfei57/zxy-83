package water

import (
	"errors"

	"trainwash/internal/store"
)

var (
	ErrNotDrainable = errors.New("water: system not drainable")
	ErrDrainBlocked = errors.New("water: drain blocked by high detergent dose")
)

func (s *System) Drain() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state != StateRinseDone && s.state != StateStopped {
		return ErrNotDrainable
	}
	if s.drainGate != nil && !s.limits.DoseOK(s.drainGate.DoseML()) {
		return ErrDrainBlocked
	}
	s.state = StateIdle
	return store.SaveJSON(s.store, store.KeyWaterState, s.state)
}

func (s *System) restoreState() {
	var state State
	if err := store.LoadJSON(s.store, store.KeyWaterState, &state); err != nil {
		return
	}
	s.state = state
}
