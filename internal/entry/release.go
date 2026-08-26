package entry

import (
	"errors"

	"trainwash/internal/store"
)

var ErrLatched = errors.New("entry: entry gate is latched")

type EntryState struct {
	Latched bool
	WashSeq int
}

func (s *Service) Latch() error {
	if s.gateLatched {
		return ErrLatched
	}
	s.gateLatched = true
	s.washSeq++
	return store.SaveJSON(s.store, store.KeyEntryLatch, EntryState{Latched: true, WashSeq: s.washSeq})
}

func (s *Service) ReleaseLatch() error {
	s.gateLatched = false
	return store.SaveJSON(s.store, store.KeyEntryLatch, EntryState{Latched: false, WashSeq: s.washSeq})
}

func (s *Service) Latched() bool {
	return s.gateLatched
}

func (s *Service) WashSeq() int {
	return s.washSeq
}

func (s *Service) restoreLatch() {
	var state EntryState
	if err := store.LoadJSON(s.store, store.KeyEntryLatch, &state); err != nil {
		return
	}
	s.gateLatched = state.Latched
	s.washSeq = state.WashSeq
}
