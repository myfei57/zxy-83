package dry

func (s *Service) FanSpeed() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.running {
		return 1800
	}
	return 0
}

func (s *Service) DutyCycle() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.running {
		return 0.92
	}
	return 0.0
}
