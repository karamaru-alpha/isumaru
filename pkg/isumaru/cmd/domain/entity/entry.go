package entity

import (
	"time"
)

type Entry struct {
	ID        string
	Type      EntryType
	Time      time.Time
	TargetIDs []string
}

type Entries []*Entry

type EntryType int32

const (
	EntryTypeMysql EntryType = iota + 1
)
