package water

import "trainwash/internal/store"

func (s *System) BeginRinse() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state != StateFlowing {
		return ErrNotFlowing
	}
	s.state = StateRinsing
	return nil
}

func (s *System) RinseDone() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state != StateRinsing {
		return ErrNotRinsing
	}
	s.state = StateRinseDone
	return store.SaveJSON(s.store, store.KeyWaterState, s.state)
}

func (s *System) RinseComplete() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state == StateRinseDone
}
