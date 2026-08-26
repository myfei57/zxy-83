package pos

import (
	"errors"

	"trainwash/internal/store"
)

var ErrBadPosition = errors.New("pos: invalid position")

func (t *Tracker) Persist(p Position) error {
	if p.FrontMM < 0 || p.LengthMM <= 0 || p.TrainID == "" {
		return ErrBadPosition
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := store.SaveJSON(t.store, store.KeyPosPosition, p); err != nil {
		return err
	}
	t.current = p
	t.persisted = true
	return nil
}

func (t *Tracker) Restore() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.restorePositionLocked()
}

func (t *Tracker) restorePosition() {
	t.mu.Lock()
	defer t.mu.Unlock()
	_ = t.restorePositionLocked()
}

func (t *Tracker) restorePositionLocked() error {
	var current Position
	if err := store.LoadJSON(t.store, store.KeyPosPosition, &current); err != nil {
		return err
	}
	t.current = current
	t.persisted = true
	return nil
}

func (t *Tracker) ClearPosition() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := t.store.Delete(store.KeyPosPosition); err != nil {
		return err
	}
	t.current = Position{}
	t.persisted = false
	return nil
}
