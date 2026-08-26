package brush

import "errors"

var ErrBadTrainLength = errors.New("brush: invalid train length for scrubbing")

func (s *Set) Scrub(trainLengthMM int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.lowered {
		return 0, ErrNotLowered
	}
	if trainLengthMM <= 0 {
		return 0, ErrBadTrainLength
	}
	if s.active.SpeedMMS <= 0 {
		return 0, ErrNoGroup
	}
	durationMS := (trainLengthMM / s.active.SpeedMMS) * 1000
	return durationMS, nil
}

func (s *Set) ScrubSweeps(trainLengthMM, gapMM int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.lowered {
		return 0, ErrNotLowered
	}
	if trainLengthMM <= 0 || gapMM < 0 {
		return 0, ErrBadTrainLength
	}
	if gapMM == 0 {
		return 1, nil
	}
	sweeps := trainLengthMM / gapMM
	if sweeps < 1 {
		sweeps = 1
	}
	return sweeps, nil
}
