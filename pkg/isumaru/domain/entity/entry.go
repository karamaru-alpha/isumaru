package entity

import (
	"time"
)

type Entry struct {
	ID       string
	Type     EntryType
	FileSize int64
	Time     time.Time
}

type Entries []*Entry

type EntryType int32

const (
	EntryTypeMysql EntryType = iota + 1
)
