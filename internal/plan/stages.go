package plan

type Stage int

const (
	StageIdle Stage = iota
	StageWash
	StageRinse
	StageDry
	StageDone
)

func (s Stage) String() string {
	switch s {
	case StageIdle:
		return "idle"
	case StageWash:
		return "wash"
	case StageRinse:
		return "rinse"
	case StageDry:
		return "dry"
	case StageDone:
		return "done"
	default:
		return "unknown"
	}
}
