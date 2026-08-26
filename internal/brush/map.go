package brush

import "trainwash/internal/store"

func (s *Set) RefreshMap() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active.Name == "" {
		return ErrNoGroup
	}
	return store.SaveJSON(s.store, store.KeyBrushGroup, s.active)
}
