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
)

type Service struct {
	water   *water.System
	mu      sync.RWMutex
	running bool
}

func NewService(waterSystem *water.System) *Service {
	return &Service{water: waterSystem}
}

func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.water.RinseComplete() {
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

// ResetCycle clears any leftover running state from a prior wash so a fresh
// cycle cannot inherit a stale "fans on" flag. StartWash calls this before the
// new cycle begins; without it the dryer would keep reporting running across
// an interrupted cycle, letting the fans spin before the new rinse completes.
func (s *Service) ResetCycle() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}
