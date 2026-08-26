package plan

type WashConfig struct {
	SprayMS      int
	RinseLevelMM int
	DrySeconds   int
	ScrubGapMM   int
}

func DefaultWashConfig() WashConfig {
	return WashConfig{
		SprayMS:      15000,
		RinseLevelMM: 500,
		DrySeconds:   30,
		ScrubGapMM:   200,
	}
}

func (c WashConfig) Validate() error {
	if c.SprayMS <= 0 {
		return ErrBadConfig
	}
	if c.RinseLevelMM <= 0 {
		return ErrBadConfig
	}
	if c.DrySeconds <= 0 {
		return ErrBadConfig
	}
	if c.ScrubGapMM < 0 {
		return ErrBadConfig
	}
	return nil
}
