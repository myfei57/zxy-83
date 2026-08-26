package pos

import (
	"sync"

	"trainwash/internal/store"
)

type Position struct {
	TrainID  string
	FrontMM  int
	LengthMM int
	ZeroMM   int
}

type Tracker struct {
	store      *store.FileStore
	mu         sync.RWMutex
	zeroMM     int
	current    Position
	persisted  bool
	headSeen   bool
	headRecord HeadRecord
	headEpoch  int
}

func NewTracker(st *store.FileStore) *Tracker {
	t := &Tracker{store: st}
	t.zeroMM = t.loadZero()
	t.restorePosition()
	t.restoreHead()
	return t
}

func (t *Tracker) ZeroMM() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.zeroMM
}

func (t *Tracker) Persisted() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.persisted
}

func (t *Tracker) HeadArrived() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.headSeen
}

func (t *Tracker) HeadMM() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.headRecord.HeadMM
}

func (t *Tracker) Current() Position {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.current
}
