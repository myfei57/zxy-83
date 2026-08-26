package brush

func (s *Set) Retract() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.lowered {
		return ErrNotLowered
	}
	s.lowered = false
	return s.resetLatchLocked()
}
