package audit

import (
	"time"

	"github.com/google/uuid"

	"trainwash/internal/store"
)

type Entry struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	TrainID   string    `json:"train_id"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

type Recorder struct {
	store *store.FileStore
}

func NewRecorder(st *store.FileStore) *Recorder {
	return &Recorder{store: st}
}

func (r *Recorder) Add(action, trainID, detail string) (Entry, error) {
	entry := Entry{
		ID:        uuid.NewString(),
		Action:    action,
		TrainID:   trainID,
		Detail:    detail,
		CreatedAt: time.Now().UTC(),
	}
	records, err := r.load()
	if err != nil {
		return Entry{}, err
	}
	records = append(records, entry)
	if len(records) > 500 {
		records = records[len(records)-500:]
	}
	if err := store.SaveJSON(r.store, store.KeyAudit, records); err != nil {
		return Entry{}, err
	}
	return entry, nil
}
