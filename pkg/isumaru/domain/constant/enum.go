package constant

type EntryTargetStatusType int32

const (
	EntryTargetStatusTypeProgress = iota + 1
	EntryTargetStatusTypeSuccess
	EntryTargetStatusTypeFailure
)

func (t EntryTargetStatusType) String() string {
	switch t {
	case EntryTargetStatusTypeProgress:
		return "EntryTargetStatusTypeProgress"
	case EntryTargetStatusTypeSuccess:
		return "EntryTargetStatusTypeSuccess"
	case EntryTargetStatusTypeFailure:
		return "EntryTargetStatusTypeFailure"
	default:
		return "Unknown"
	}
}

type TargetType int32

const (
	TargetTypeSlowQueryLog TargetType = iota + 1
	TargetTypeAccessLog
	TargetTypeFgprof
)

func (t TargetType) String() string {
	switch t {
	case TargetTypeSlowQueryLog:
		return "SlowQueryLog"
	case TargetTypeAccessLog:
		return "AccessLog"
	case TargetTypeFgprof:
		return "Fgprof"
	default:
		return "Unknown"
	}
}
