package ns

type SegmentKind int

const (
	SegmentHead SegmentKind = iota
	SegmentBody
	SegmentTail
)

func (k SegmentKind) String() string {
	switch k {
	case SegmentHead:
		return "head"
	case SegmentBody:
		return "body"
	case SegmentTail:
		return "tail"
	default:
		return "unknown"
	}
}

type Segment struct {
	Kind    SegmentKind
	StartMM int
	EndMM   int
}

func (s Segment) Contains(mm int) bool {
	return mm >= s.StartMM && mm < s.EndMM
}

func (s Segment) LengthMM() int {
	if s.EndMM < s.StartMM {
		return 0
	}
	return s.EndMM - s.StartMM
}

func SegmentsOf(lengthMM int) []Segment {
	if lengthMM <= 0 {
		lengthMM = 1
	}
	head := lengthMM / 4
	tail := lengthMM / 5
	if head < 500 {
		head = 500
	}
	if tail < 400 {
		tail = 400
	}
	if head+tail >= lengthMM {
		tail = lengthMM / 10
	}
	bodyEnd := lengthMM - tail
	if bodyEnd < head {
		bodyEnd = head
	}
	return []Segment{
		{Kind: SegmentHead, StartMM: 0, EndMM: head},
		{Kind: SegmentBody, StartMM: head, EndMM: bodyEnd},
		{Kind: SegmentTail, StartMM: bodyEnd, EndMM: lengthMM},
	}
}

type ZoneKind int

const (
	ZoneSide ZoneKind = iota
	ZoneRoof
)

func (k ZoneKind) String() string {
	switch k {
	case ZoneSide:
		return "side"
	case ZoneRoof:
		return "roof"
	default:
		return "unknown"
	}
}

type Zone struct {
	Kind    ZoneKind
	StartMM int
	EndMM   int
}

func (z Zone) Contains(mm int) bool {
	return mm >= z.StartMM && mm < z.EndMM
}

type TrainProfile struct {
	LengthMM int
	Zones    []Zone
}

func ProfileOf(lengthMM int, roofCovered bool) TrainProfile {
	zones := []Zone{
		{Kind: ZoneSide, StartMM: 0, EndMM: lengthMM},
	}
	if roofCovered {
		roofStart := lengthMM / 5
		roofEnd := lengthMM - roofStart
		if roofEnd <= roofStart {
			roofEnd = roofStart + 1
		}
		zones = append(zones, Zone{Kind: ZoneRoof, StartMM: roofStart, EndMM: roofEnd})
	}
	return TrainProfile{LengthMM: lengthMM, Zones: zones}
}
