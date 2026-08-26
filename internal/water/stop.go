package water

import "trainwash/internal/store"

func (s *System) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state == StateIdle || s.state == StateStopped {
		return nil
	}
	s.state = StateStopped
	return store.SaveJSON(s.store, store.KeyWaterState, s.state)
}

func (s *System) StopReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state == StateStopped || s.state == StateIdle
}
