package ns

import "fmt"

type StationKind int

const (
	StationInbound StationKind = iota
	StationWash
	StationRinse
	StationDry
	StationOutbound
)

func (k StationKind) String() string {
	switch k {
	case StationInbound:
		return "inbound"
	case StationWash:
		return "wash"
	case StationRinse:
		return "rinse"
	case StationDry:
		return "dry"
	case StationOutbound:
		return "outbound"
	default:
		return "unknown"
	}
}

type Station struct {
	ID       string
	Index    int
	OffsetMM int
	Kind     StationKind
}

type StationLayout struct {
	stations []Station
}

func NewStationLayout() StationLayout {
	return StationLayout{stations: []Station{
		{ID: "in-01", Index: 0, OffsetMM: 0, Kind: StationInbound},
		{ID: "wash-01", Index: 1, OffsetMM: 8000, Kind: StationWash},
		{ID: "rinse-01", Index: 2, OffsetMM: 16000, Kind: StationRinse},
		{ID: "dry-01", Index: 3, OffsetMM: 24000, Kind: StationDry},
		{ID: "out-01", Index: 4, OffsetMM: 32000, Kind: StationOutbound},
	}}
}

func (l StationLayout) Count() int {
	return len(l.stations)
}

func (l StationLayout) StationAt(mm int) (Station, bool) {
	best := -1
	for i, station := range l.stations {
		if mm >= station.OffsetMM {
			best = i
		}
	}
	if best < 0 {
		return Station{}, false
	}
	return l.stations[best], true
}

func (l StationLayout) StationByKind(kind StationKind) (Station, bool) {
	for _, station := range l.stations {
		if station.Kind == kind {
			return station, true
		}
	}
	return Station{}, false
}

func (l StationLayout) NearestWashOffset(frontMM int) int {
	wash, ok := l.StationByKind(StationWash)
	if !ok {
		return 0
	}
	return wash.OffsetMM - frontMM
}

func (l StationLayout) String() string {
	names := make([]string, 0, len(l.stations))
	for _, station := range l.stations {
		names = append(names, fmt.Sprintf("%s@%d", station.ID, station.OffsetMM))
	}
	return fmt.Sprint(names)
}
