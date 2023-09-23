package entity

import (
	"time"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
)

type Target struct {
	ID       string
	Type     constant.TargetType
	URL      string
	Path     string
	Duration time.Duration
}

type Targets []*Target
