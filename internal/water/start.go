package water

import (
	"errors"
	"sync"

	"trainwash/internal/ns"
	"trainwash/internal/store"
)

type State int

const (
	StateIdle State = iota
	StateFlowing
	StateRinsing
	StateRinseDone
	StateStopped
)

func (st State) String() string {
	switch st {
	case StateIdle:
		return "idle"
	case StateFlowing:
		return "flowing"
	case StateRinsing:
		return "rinsing"
	case StateRinseDone:
		return "rinse_done"
	case StateStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

var (
	ErrAlreadyRunning = errors.New("water: system already running")
	ErrNotFlowing     = errors.New("water: system is not flowing")
	ErrNotRinsing     = errors.New("water: rinse has not started")
	ErrBadLevel       = errors.New("water: invalid pressure level")
	ErrPressureRange  = errors.New("water: pressure out of range")
	ErrStaleBaseline  = errors.New("water: pressure baseline is stale")
)

type DrainGate interface {
	DoseML() float64
}

type System struct {
	store        *store.FileStore
	limits       ns.Limits
	mu           sync.RWMutex
	state        State
	gainMPA      float64
	gainRevision int
	cycleEpoch   int
	drainGate    DrainGate
}

func NewSystem(st *store.FileStore, limits ns.Limits) *System {
	s := &System{store: st, limits: limits, gainMPA: 5.0}
	s.restoreGain()
	s.restoreState()
	s.restoreCycle()
	return s
}

func (s *System) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state == StateFlowing || s.state == StateRinsing || s.state == StateRinseDone {
		return ErrAlreadyRunning
	}
	s.state = StateFlowing
	s.cycleEpoch++
	if err := store.SaveJSON(s.store, store.KeyWaterState, s.state); err != nil {
		return err
	}
	return store.SaveJSON(s.store, store.KeyWaterCycle, s.cycleEpoch)
}

func (s *System) State() State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

func (s *System) IsFlowing() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state == StateFlowing
}

func (s *System) IsStopped() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state == StateStopped
}

func (s *System) AttachDrainGate(gate DrainGate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.drainGate = gate
}

func (s *System) CycleID() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cycleEpoch
}

func (s *System) restoreCycle() {
	var cycle int
	if err := store.LoadJSON(s.store, store.KeyWaterCycle, &cycle); err != nil {
		return
	}
	s.cycleEpoch = cycle
}
