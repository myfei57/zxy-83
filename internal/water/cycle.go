package water

func (s *System) CanStart() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state == StateIdle || s.state == StateStopped
}

func (s *System) CanRinse() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state == StateFlowing
}

func (s *System) CanDrain() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state == StateRinseDone || s.state == StateStopped
}
