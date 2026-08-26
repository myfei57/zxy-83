package audit

import (
	"path/filepath"
	"testing"

	"trainwash/internal/store"
)

func TestRecorderAddAndList(t *testing.T) {
	recorder := NewRecorder(store.New(filepath.Join(t.TempDir(), "audit")))
	entry, err := recorder.Add("wash_start", "T-1", "cycle started")
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if entry.ID == "" {
		t.Fatal("entry should have a uuid id")
	}
	if _, err := recorder.Add("wash_complete", "T-1", "cycle done"); err != nil {
		t.Fatalf("add second: %v", err)
	}
	records, err := recorder.List(10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	count, err := recorder.Count()
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}
