package pos

import (
	"errors"

	"trainwash/internal/store"
)

var ErrBadHead = errors.New("pos: invalid head position")

type HeadRecord struct {
	HeadMM int
	Epoch  int
}

func (t *Tracker) MarkHead(headMM int) error {
	if headMM < 0 {
		return ErrBadHead
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.headEpoch++
	record := HeadRecord{HeadMM: headMM, Epoch: t.headEpoch}
	if err := store.SaveJSON(t.store, store.KeyPosHead, record); err != nil {
		return err
	}
	return nil
}

func (t *Tracker) ClearHead() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.headEpoch++
	record := HeadRecord{HeadMM: 0, Epoch: t.headEpoch}
	if err := store.SaveJSON(t.store, store.KeyPosHead, record); err != nil {
		return err
	}
	return nil
}

func (t *Tracker) HeadRecord() (HeadRecord, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.headRecord, t.headSeen
}

func (t *Tracker) HeadEpoch() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.headEpoch
}

func (t *Tracker) restoreHead() {
	t.mu.Lock()
	defer t.mu.Unlock()
}
