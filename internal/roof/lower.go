package roof

import (
	"errors"
	"sync"

	"trainwash/internal/pos"
)

var (
	ErrAlreadyLowered = errors.New("roof: roof brush already lowered")
	ErrNotLowered     = errors.New("roof: roof brush not lowered")
	ErrHeadNotReady   = errors.New("roof: train head has not reached the work position")
)

type Service struct {
	pos     *pos.Tracker
	mu      sync.RWMutex
	lowered bool
}

func NewService(tracker *pos.Tracker) *Service {
	return &Service{pos: tracker}
}

// Lower drops the roof brush onto the train roof. The train head must have
// arrived at the work position first (see HeadReady); otherwise the brush
// would land on the windshield before the roof surface is under it. This guard
// must live here rather than only at the console layer so that every caller —
// manual console requests, the wash cycle, and any future automation — is
// protected from dropping the brush before the head is in position.
func (s *Service) Lower() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lowered {
		return ErrAlreadyLowered
	}
	if !s.pos.HeadArrived() {
		return ErrHeadNotReady
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
