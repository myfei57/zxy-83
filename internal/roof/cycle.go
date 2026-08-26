package roof

func (s *Service) Brush(segmentStartMM, segmentEndMM int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.lowered {
		return 0, ErrNotLowered
	}
	start, end := s.pos.SegmentPoint(segmentStartMM, segmentEndMM)
	return end - start, nil
}

func (s *Service) BrushZone(offsetMM int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.lowered {
		return 0, ErrNotLowered
	}
	return s.pos.PointFor(offsetMM), nil
}
