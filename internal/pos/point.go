package pos

func (t *Tracker) PointFor(zoneOffsetMM int) int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.zeroMM + t.current.FrontMM + zoneOffsetMM
}

func (t *Tracker) SegmentPoint(segmentStartMM, segmentEndMM int) (int, int) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	start := t.zeroMM + t.current.FrontMM + segmentStartMM
	end := t.zeroMM + t.current.FrontMM + segmentEndMM
	return start, end
}
