package pos

import (
	"errors"

	"trainwash/internal/store"
)

var ErrBadZero = errors.New("pos: invalid zero")

func (t *Tracker) Recalibrate(zeroMM int) error {
	if zeroMM < 0 || zeroMM > 5000 {
		return ErrBadZero
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := store.SaveJSON(t.store, store.KeyPosZero, zeroMM); err != nil {
		return err
	}
	t.zeroMM = zeroMM
	return nil
}

func (t *Tracker) RefreshZero(zeroMM int) error {
	return t.Recalibrate(zeroMM)
}

func (t *Tracker) loadZero() int {
	var zero int
	if err := store.LoadJSON(t.store, store.KeyPosZero, &zero); err != nil {
		return 0
	}
	return zero
}
