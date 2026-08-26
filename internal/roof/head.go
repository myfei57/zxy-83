package roof

func (s *Service) HeadReady() bool {
	return s.pos.HeadArrived()
}

func (s *Service) HeadPositionMM() int {
	return s.pos.HeadMM()
}
