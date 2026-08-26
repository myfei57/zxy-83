package brush

import (
	"sync"

	"trainwash/internal/ns"
	"trainwash/internal/pos"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

type Group struct {
	Name     string
	LengthMM int
	SpeedMMS int
	Zones    []ns.Zone
}

func GroupFromSpec(name string, lengthMM, speedMMS int) Group {
	profile := ns.ProfileOf(lengthMM, true)
	return Group{
		Name:     name,
		LengthMM: profile.LengthMM,
		SpeedMMS: speedMMS,
		Zones:    profile.Zones,
	}
}

type Set struct {
	store           *store.FileStore
	pos             *pos.Tracker
	water           *water.System
	releaser        LatchReleaser
	mu              sync.RWMutex
	active          Group
	lowered         bool
	zeroMM          int
	gainRevision    int
	lastPressureMPA float64
}

func NewSet(st *store.FileStore, tracker *pos.Tracker, waterSystem *water.System) *Set {
	s := &Set{store: st, pos: tracker, water: waterSystem}
	s.gainRevision = waterSystem.GainRevision()
	s.restoreFromStore()
	return s
}

func (s *Set) ApplyGroup(group Group) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if group.Name == "" || group.LengthMM <= 0 || group.SpeedMMS <= 0 {
		return ErrNoGroup
	}
	s.active = group
	s.zeroMM = s.pos.ZeroMM()
	return store.SaveJSON(s.store, store.KeyBrushGroup, group)
}

func (s *Set) ActiveGroup() Group {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.active
}

func (s *Set) PublishGroup(name string, lengthMM, speedMMS int) error {
	group := GroupFromSpec(name, lengthMM, speedMMS)
	s.mu.Lock()
	defer s.mu.Unlock()
	if group.Name == "" || group.LengthMM <= 0 || group.SpeedMMS <= 0 {
		return ErrNoGroup
	}
	return store.SaveJSON(s.store, store.KeyBrushGroup, group)
}

func (s *Set) restoreFromStore() {
	var group Group
	if err := store.LoadJSON(s.store, store.KeyBrushGroup, &group); err != nil {
		return
	}
	s.active = group
}
