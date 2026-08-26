package roof

import (
	"errors"
	"sync"

	"trainwash/internal/pos"
)

var (
	ErrHeadNotArrived = errors.New("roof: train head has not arrived")
	ErrAlreadyLowered = errors.New("roof: roof brush already lowered")
	ErrNotLowered     = errors.New("roof: roof brush not lowered")
)

type Service struct {
	pos           *pos.Tracker
	mu            sync.RWMutex
	lowered       bool
	expectedEpoch int
}

func NewService(tracker *pos.Tracker) *Service {
	s := &Service{pos: tracker}
	s.syncEpochLocked()
	return s
}

func (s *Service) Lower() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.pos.HeadRecord()
	if ok && record.Epoch > s.expectedEpoch {
		s.expectedEpoch = record.Epoch
	}
	if !ok || record.HeadMM <= 0 || record.Epoch < s.expectedEpoch {
		return ErrHeadNotArrived
	}
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

func (s *Service) syncEpochLocked() {
	s.expectedEpoch = s.pos.HeadEpoch()
}
