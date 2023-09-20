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
	TargetTypeAccessLog TargetType = iota + 1
	TargetTypeSlowQueryLog
)

func (t TargetType) String() string {
	switch t {
	case TargetTypeAccessLog:
		return "AccessLog"
	case TargetTypeSlowQueryLog:
		return "SlowQueryLog"
	default:
		return ""
	}
}
