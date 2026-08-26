package verifycase

import (
	"path/filepath"
	"testing"

	"trainwash/internal/brush"
	"trainwash/internal/entry"
	"trainwash/internal/ns"
	"trainwash/internal/pos"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

func TestTwsBrushMapFresh(t *testing.T) {
	root := filepath.Join(t.TempDir(), "state")
	st := store.New(root)
	limits := ns.DefaultLimits()
	waterSystem := water.NewSystem(st, limits)
	brushSet := brush.NewSet(st, pos.NewTracker(st), waterSystem)
	entryService := entry.NewService(st)
	entryService.AttachPublisher(brushSet)
	if err := entryService.TypeChange(entry.TypeLong); err != nil {
		t.Fatalf("type change: %v", err)
	}
	if got := brushSet.ActiveGroup().Name; got != "long-set" {
		t.Fatalf("brush group must follow the train type change, got %q", got)
	}
	restarted := entry.NewService(st)
	restarted.AttachPublisher(brushSet)
	if got := brushSet.ActiveGroup().Name; got != "long-set" {
		t.Fatalf("brush group must follow the restored train type, got %q", got)
	}
	if err := entryService.TypeChange(entry.TypeShort); err != nil {
		t.Fatalf("type change back: %v", err)
	}
	if got := brushSet.ActiveGroup().Name; got != "short-set" {
		t.Fatalf("brush group must follow the second type change, got %q", got)
	}
}
