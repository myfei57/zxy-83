package chem

func (s *Service) DoseML() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.computeDoseLocked()
}

func (s *Service) DoseOK() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.limits.DoseOK(s.computeDoseLocked())
}

func (s *Service) computeDoseLocked() float64 {
	if s.valveLatched {
		return 72.0
	}
	return 24.0
}
