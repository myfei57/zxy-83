package dry

import (
	"errors"
	"sync"

	"trainwash/internal/water"
)

var (
	ErrRinseNotDone   = errors.New("dry: rinse has not completed")
	ErrAlreadyRunning = errors.New("dry: dryer already running")
	ErrNotRunning     = errors.New("dry: dryer is not running")
	ErrBadDuration    = errors.New("dry: invalid drying duration")
	ErrCycleMismatch  = errors.New("dry: wash cycle does not match")
)

type Service struct {
	water         *water.System
	mu            sync.RWMutex
	running       bool
	expectedCycle int
}

func NewService(waterSystem *water.System) *Service {
	return &Service{water: waterSystem, expectedCycle: waterSystem.CycleID()}
}

func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.water.CycleID() != s.expectedCycle {
		return ErrCycleMismatch
	}
	if s.water.State() != water.StateRinseDone {
		return ErrRinseNotDone
	}
	if s.running {
		return ErrAlreadyRunning
	}
	s.running = true
	return nil
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return ErrNotRunning
	}
	s.running = false
	return nil
}

func (s *Service) Running() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *Service) DryFor(seconds int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return ErrNotRunning
	}
	if seconds <= 0 {
		return ErrBadDuration
	}
	return nil
}

func (s *Service) ResetCycle() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expectedCycle = s.water.CycleID()
}
