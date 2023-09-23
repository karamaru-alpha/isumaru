package entity

import (
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
)

type Entry struct {
	ID      string
	Targets EntryTargets
}
type Entries []*Entry

type EntryTarget struct {
	*Target
	StatusType constant.EntryTargetStatusType
	Error      error
}

type EntryTargets []*EntryTarget
