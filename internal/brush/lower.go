package brush

import "errors"

var (
	ErrNotPersisted   = errors.New("brush: position not persisted")
	ErrAlreadyLowered = errors.New("brush: brushes already lowered")
	ErrNotLowered     = errors.New("brush: brushes not lowered")
	ErrNoGroup        = errors.New("brush: no active brush group")
)

func (s *Set) LowerSide() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.pos.Persisted() {
		return ErrNotPersisted
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
