package chem

import "trainwash/internal/store"

type ValveState struct {
	Latched  bool
	AlarmSeq int
}

func (s *Service) ReleaseValve() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.valveLatched = false
	return store.SaveJSON(s.store, store.KeyChemValve, ValveState{Latched: false, AlarmSeq: s.alarmSeq})
}

func (s *Service) ValveLatched() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.valveLatched
}

func (s *Service) AlarmSeq() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.alarmSeq
}

func (s *Service) restoreValve() {
	var state ValveState
	if err := store.LoadJSON(s.store, store.KeyChemValve, &state); err != nil {
		return
	}
	s.valveLatched = state.Latched
	s.alarmSeq = state.AlarmSeq
}
