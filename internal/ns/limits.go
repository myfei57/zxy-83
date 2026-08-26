package ns

type Limits struct {
	MinPressureMPA float64
	MaxPressureMPA float64
	MinDoseML      float64
	MaxDoseML      float64
	MaxSpeedMMS    int
	MaxSprayMS     int
}

func DefaultLimits() Limits {
	return Limits{
		MinPressureMPA: 2.0,
		MaxPressureMPA: 9.0,
		MinDoseML:      10.0,
		MaxDoseML:      80.0,
		MaxSpeedMMS:    1200,
		MaxSprayMS:     30000,
	}
}

func (l Limits) PressureOK(mpa float64) bool {
	return mpa >= l.MinPressureMPA && mpa <= l.MaxPressureMPA
}

func (l Limits) DoseOK(ml float64) bool {
	return ml >= l.MinDoseML && ml <= l.MaxDoseML
}

func (l Limits) SpeedOK(mmPerSecond int) bool {
	return mmPerSecond > 0 && mmPerSecond <= l.MaxSpeedMMS
}

func (l Limits) SprayDurationOK(ms int) bool {
	return ms > 0 && ms <= l.MaxSprayMS
}
