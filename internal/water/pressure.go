package water

func (s *System) PressureForRevision(revision, levelMM int) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if levelMM <= 0 {
		return 0, ErrBadLevel
	}
	if revision != s.gainRevision {
		return 0, ErrStaleBaseline
	}
	mpa := s.gainMPA * float64(levelMM) / 1000.0
	if !s.limits.PressureOK(mpa) {
		return 0, ErrPressureRange
	}
	return mpa, nil
}

func (s *System) ExpectedPressure(gainMPA float64, levelMM int) (float64, error) {
	if gainMPA <= 0 || levelMM <= 0 {
		return 0, ErrBadGain
	}
	mpa := gainMPA * float64(levelMM) / 1000.0
	if !s.limits.PressureOK(mpa) {
		return 0, ErrPressureRange
	}
	return mpa, nil
}
