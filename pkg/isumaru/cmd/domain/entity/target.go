package entity

import (
	"time"
)

type Target struct {
	ID       string
	Type     TargetType
	URL      string
	Path     string
	Duration time.Duration
}

type Targets []*Target

type TargetType int32

const (
	TargetTypeSlowQueryLog TargetType = iota + 1
	TargetTypeAccessLog
)

func (t TargetType) String() string {
	switch t {
	case TargetTypeSlowQueryLog:
		return "SlowQueryLog"
	case TargetTypeAccessLog:
		return "AccessLog"
	default:
		return ""
	}
}
