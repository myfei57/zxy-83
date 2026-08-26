package brush

import "errors"

var (
	ErrBadRinseLevel = errors.New("brush: invalid rinse level")
)

func (s *Set) Rinse(levelMM int) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.lowered {
		return 0, ErrNotLowered
	}
	if levelMM <= 0 {
		return 0, ErrBadRinseLevel
	}
	s.refreshGainLocked()
	mpa, err := s.water.PressureForRevision(s.gainRevision, levelMM)
	if err != nil {
		return 0, err
	}
	s.lastPressureMPA = mpa
	return mpa, nil
}

func (s *Set) RefreshGain() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refreshGainLocked()
}

func (s *Set) refreshGainLocked() {
	if _, ok := s.water.GainAtRevision(s.gainRevision); !ok {
		s.gainRevision = s.water.GainRevision()
	}
}

func (s *Set) LastPressureMPA() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastPressureMPA
}
