package chem

import (
	"errors"
	"sync"

	"trainwash/internal/ns"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

var (
	ErrWaterNotFlowing = errors.New("chem: water is not flowing")
	ErrBadSpray        = errors.New("chem: invalid spray duration")
)

type Service struct {
	store        *store.FileStore
	water        *water.System
	limits       ns.Limits
	mu           sync.RWMutex
	valveLatched bool
	alarm        bool
	alarmSeq     int
	lastDoseML   float64
}

func NewService(st *store.FileStore, waterSystem *water.System, limits ns.Limits) *Service {
	s := &Service{store: st, water: waterSystem, limits: limits}
	s.restoreValve()
	return s
}

func (s *Service) Spray(ms int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.water.IsFlowing() {
		return ErrWaterNotFlowing
	}
	if !s.limits.SprayDurationOK(ms) {
		return ErrBadSpray
	}
	s.lastDoseML = s.computeDoseLocked()
	return nil
}

func (s *Service) LastDoseML() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastDoseML
}
