package brush

import "errors"

var ErrNoReleaser = errors.New("brush: latch releaser not attached")

type LatchReleaser interface {
	ReleaseLatch() error
}

func (s *Set) AttachReleaser(releaser LatchReleaser) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.releaser = releaser
}

func (s *Set) ResetLatch() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.resetLatchLocked()
}

func (s *Set) resetLatchLocked() error {
	if s.releaser == nil {
		return ErrNoReleaser
	}
	return s.releaser.ReleaseLatch()
}
