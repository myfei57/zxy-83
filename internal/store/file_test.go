package store

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestFileStoreRoundTrip(t *testing.T) {
	root := filepath.Join(t.TempDir(), "state")
	st := New(root)
	payload := []byte(`{"k":1}`)
	if err := st.Save("pos/position", payload); err != nil {
		t.Fatalf("save: %v", err)
	}
	if !st.Exists("pos/position") {
		t.Fatal("key should exist after save")
	}
	got, err := st.Load("pos/position")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("payload mismatch: %s", got)
	}
	if err := st.Delete("pos/position"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := st.Load("pos/position"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestJSONHelpers(t *testing.T) {
	root := filepath.Join(t.TempDir(), "state")
	st := New(root)
	type sample struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	value := sample{Name: "wash", Count: 3}
	if err := SaveJSON(st, "audit/records", value); err != nil {
		t.Fatalf("save json: %v", err)
	}
	var loaded sample
	if err := LoadJSON(st, "audit/records", &loaded); err != nil {
		t.Fatalf("load json: %v", err)
	}
	if loaded != value {
		t.Fatalf("round trip mismatch: %+v", loaded)
	}
}
