package entry

type GroupSpec struct {
	Name     string
	LengthMM int
	SpeedMMS int
}

type GroupMap struct {
	Short GroupSpec
	Long  GroupSpec
}

func DefaultGroupMap() GroupMap {
	return GroupMap{
		Short: GroupSpec{Name: "short-set", LengthMM: 20000, SpeedMMS: 600},
		Long:  GroupSpec{Name: "long-set", LengthMM: 40000, SpeedMMS: 400},
	}
}

func (m GroupMap) Resolve(trainType TrainType) GroupSpec {
	if trainType == TypeLong {
		return m.Long
	}
	return m.Short
}
