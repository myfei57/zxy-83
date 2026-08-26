package water

import (
	"errors"

	"trainwash/internal/store"
)

var ErrBadGain = errors.New("water: invalid pressure gain")

type GainBaseline struct {
	GainMPA  float64
	Revision int
}

func (s *System) RecalibratePressure(gainMPA float64) error {
	if gainMPA <= 0 || gainMPA > 10 {
		return ErrBadGain
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gainRevision++
	baseline := GainBaseline{GainMPA: gainMPA, Revision: s.gainRevision}
	if err := store.SaveJSON(s.store, store.KeyWaterGain, baseline); err != nil {
		return err
	}
	s.gainMPA = gainMPA
	return nil
}

func (s *System) GainMPA() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.gainMPA
}

func (s *System) GainRevision() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.gainRevision
}

func (s *System) GainAtRevision(revision int) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if revision != s.gainRevision {
		return 0, false
	}
	return s.gainMPA, true
}

func (s *System) restoreGain() {
	var baseline GainBaseline
	if err := store.LoadJSON(s.store, store.KeyWaterGain, &baseline); err != nil {
		return
	}
	s.gainMPA = baseline.GainMPA
	s.gainRevision = baseline.Revision
}
