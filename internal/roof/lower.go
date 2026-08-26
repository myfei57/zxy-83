package roof

import (
	"errors"
	"sync"

	"trainwash/internal/pos"
)

var (
	ErrAlreadyLowered = errors.New("roof: roof brush already lowered")
	ErrNotLowered     = errors.New("roof: roof brush not lowered")
)

type Service struct {
	pos     *pos.Tracker
	mu      sync.RWMutex
	lowered bool
}

func NewService(tracker *pos.Tracker) *Service {
	return &Service{pos: tracker}
}

func (s *Service) Lower() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lowered {
		return ErrAlreadyLowered
	}
	s.lowered = true
	return nil
}

func (s *Service) IsLowered() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lowered
}

func (s *Service) Raise() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.lowered {
		return ErrNotLowered
	}
	s.lowered = false
	return nil
}
