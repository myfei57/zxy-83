package entry

type TrainType string

const (
	TypeShort TrainType = "short"
	TypeLong  TrainType = "long"
)

type Train struct {
	ID       string
	Type     TrainType
	LengthMM int
}

func NewTrain(id string, trainType TrainType, lengthMM int) Train {
	return Train{ID: id, Type: NormalizeType(string(trainType)), LengthMM: lengthMM}
}

func NormalizeType(value string) TrainType {
	if value == "long" || value == "Long" || value == "LONG" {
		return TypeLong
	}
	return TypeShort
}
