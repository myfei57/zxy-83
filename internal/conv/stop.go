package conv

import (
	"errors"
	"sync"

	"trainwash/internal/ns"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

var ErrWaterRunning = errors.New("conv: water system must stop before conveyor")

type ConvState struct {
	Running    bool
	PositionMM int
}

type Service struct {
	store      *store.FileStore
	water      *water.System
	layout     ns.StationLayout
	limits     ns.Limits
	mu         sync.RWMutex
	positionMM int
	running    bool
}

func NewService(st *store.FileStore, waterSystem *water.System, layout ns.StationLayout, limits ns.Limits) *Service {
	s := &Service{store: st, water: waterSystem, layout: layout, limits: limits}
	s.restore()
	return s
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.water.IsStopped() {
		return ErrWaterRunning
	}
	s.running = false
	return store.SaveJSON(s.store, store.KeyConvState, ConvState{Running: false, PositionMM: s.positionMM})
}

func (s *Service) Running() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *Service) PositionMM() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.positionMM
}

func (s *Service) restore() {
	var state ConvState
	if err := store.LoadJSON(s.store, store.KeyConvState, &state); err != nil {
		return
	}
	s.running = state.Running
	s.positionMM = state.PositionMM
}
