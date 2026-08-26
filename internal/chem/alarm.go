package chem

import "trainwash/internal/store"

func (s *Service) SetAlarm() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alarm = true
	s.alarmSeq++
	s.valveLatched = true
	return store.SaveJSON(s.store, store.KeyChemValve, ValveState{Latched: true, AlarmSeq: s.alarmSeq})
}

func (s *Service) AlarmClear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alarm = false
	return store.SaveJSON(s.store, store.KeyChemValve, ValveState{Latched: s.valveLatched, AlarmSeq: s.alarmSeq})
}

func (s *Service) AlarmActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.alarm
}
