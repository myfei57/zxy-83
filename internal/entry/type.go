package entry

import (
	"errors"

	"trainwash/internal/store"
)

var ErrNoTrain = errors.New("entry: no train recognized")

type MapPublisher interface {
	PublishGroup(name string, lengthMM, speedMMS int) error
}

type Service struct {
	store       *store.FileStore
	publisher   MapPublisher
	current     TrainType
	mapping     GroupMap
	gateLatched bool
	washSeq     int
}

func NewService(st *store.FileStore) *Service {
	s := &Service{store: st, current: TypeShort, mapping: DefaultGroupMap()}
	s.restore()
	s.restoreLatch()
	return s
}

func (s *Service) AttachPublisher(publisher MapPublisher) {
	s.publisher = publisher
	s.publishCurrent()
}

func (s *Service) Recognize(train Train) (TrainType, error) {
	if train.ID == "" || train.LengthMM <= 0 {
		return "", ErrNoTrain
	}
	s.current = NormalizeType(string(train.Type))
	if err := store.SaveJSON(s.store, store.KeyEntryType, s.current); err != nil {
		return "", err
	}
	if err := s.publishCurrent(); err != nil {
		return "", err
	}
	return s.current, nil
}

func (s *Service) TypeChange(trainType TrainType) error {
	if trainType != TypeShort && trainType != TypeLong {
		return errors.New("entry: unsupported train type")
	}
	s.current = trainType
	if err := store.SaveJSON(s.store, store.KeyEntryType, s.current); err != nil {
		return err
	}
	return s.publishCurrent()
}

func (s *Service) CurrentType() TrainType {
	return s.current
}

func (s *Service) ResolveGroup(trainType TrainType) GroupSpec {
	return s.mapping.Resolve(trainType)
}

func (s *Service) publishCurrent() error {
	if s.publisher == nil {
		return nil
	}
	spec := s.mapping.Resolve(s.current)
	return s.publisher.PublishGroup(spec.Name, spec.LengthMM, spec.SpeedMMS)
}

func (s *Service) restore() {
	var current TrainType
	if err := store.LoadJSON(s.store, store.KeyEntryType, &current); err != nil {
		current = TypeShort
	}
	s.current = current
}
