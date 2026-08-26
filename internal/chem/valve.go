package chem

type ValveState struct {
	Latched  bool
	AlarmSeq int
}

func (s *Service) ReleaseValve() error {
	return nil
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
}
