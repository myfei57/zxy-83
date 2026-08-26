package audit

import (
	"errors"

	"trainwash/internal/store"
)

func (r *Recorder) List(limit int) ([]Entry, error) {
	records, err := r.load()
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > len(records) {
		limit = len(records)
	}
	return records[len(records)-limit:], nil
}

func (r *Recorder) Count() (int, error) {
	records, err := r.load()
	if err != nil {
		return 0, err
	}
	return len(records), nil
}

func (r *Recorder) load() ([]Entry, error) {
	var records []Entry
	err := store.LoadJSON(r.store, store.KeyAudit, &records)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return records, nil
}
