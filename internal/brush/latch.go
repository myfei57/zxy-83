package brush

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
	return nil
}
