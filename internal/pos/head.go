package pos

import (
	"errors"

	"trainwash/internal/store"
)

var ErrBadHead = errors.New("pos: invalid head position")

// HeadRecord captures the last reported train-head position. Arrived records
// whether the head has reached the work position; it is the source of truth for
// HeadArrived and must be persisted so the state survives a restart.
type HeadRecord struct {
	HeadMM  int
	Epoch   int
	Arrived bool
}

func (t *Tracker) MarkHead(headMM int) error {
	if headMM < 0 {
		return ErrBadHead
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.headEpoch++
	record := HeadRecord{HeadMM: headMM, Epoch: t.headEpoch, Arrived: true}
	if err := store.SaveJSON(t.store, store.KeyPosHead, record); err != nil {
		return err
	}
	t.headRecord = record
	t.headSeen = true
	return nil
}

func (t *Tracker) ClearHead() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.headEpoch++
	record := HeadRecord{HeadMM: 0, Epoch: t.headEpoch, Arrived: false}
	if err := store.SaveJSON(t.store, store.KeyPosHead, record); err != nil {
		return err
	}
	t.headRecord = record
	t.headSeen = false
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

// restoreHead reloads the last persisted head record so that the head-arrival
// state survives a restart. Without it the in-memory headSeen flag stays false
// even when a head was already marked on disk.
func (t *Tracker) restoreHead() {
	t.mu.Lock()
	defer t.mu.Unlock()
	var record HeadRecord
	if err := store.LoadJSON(t.store, store.KeyPosHead, &record); err != nil {
		return
	}
	t.headRecord = record
	t.headEpoch = record.Epoch
	t.headSeen = record.Arrived
}
