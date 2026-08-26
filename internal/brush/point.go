package brush

func (s *Set) RefreshZero() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.zeroMM = s.pos.ZeroMM()
}

func (s *Set) PointFor(offsetMM int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.pos.PointFor(offsetMM)
}

func (s *Set) CachedZeroMM() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.zeroMM
}
