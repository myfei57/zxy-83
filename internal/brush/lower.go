package brush

import "errors"

var (
	ErrNotPersisted   = errors.New("brush: position not persisted")
	ErrAlreadyLowered = errors.New("brush: brushes already lowered")
	ErrNotLowered     = errors.New("brush: brushes not lowered")
	ErrNoGroup        = errors.New("brush: no active brush group")
	ErrHeadNotReady   = errors.New("brush: train head has not reached the work position")
)

// LowerSide drops the side brushes onto the train. The train head must have
// arrived at the work position first; otherwise the brushes land on the
// windshield before the body is under them. This guard must live here rather
// than only at the console layer so that every caller — manual console
// requests, the wash cycle, and any future automation — is protected from
// dropping the brushes before the head is in position.
func (s *Set) LowerSide() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.pos.Persisted() {
		return ErrNotPersisted
	}
	if !s.pos.HeadArrived() {
		return ErrHeadNotReady
	}
	if s.lowered {
		return ErrAlreadyLowered
	}
	if s.active.Name == "" {
		return ErrNoGroup
	}
	s.zeroMM = s.pos.ZeroMM()
	s.lowered = true
	return nil
}

func (s *Set) IsLowered() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lowered
}

func (s *Set) RaiseSide() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.lowered {
		return ErrNotLowered
	}
	s.lowered = false
	return nil
}
