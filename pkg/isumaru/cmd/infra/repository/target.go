package repository

import (
	"context"
	"time"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/repository"
)

type targetRepository struct{}

func NewTargetRepository() repository.TargetRepository {
	return &targetRepository{}
}

var targets = entity.Targets{
	{
		ID:       "isu1",
		Duration: time.Second * 70,
		Type:     entity.TargetTypeAccessLog,
		URL:      "http://localhost:19000",
		Path:     constant.DefaultAccessLogPath,
	},
	{
		ID:       "isu1",
		Duration: time.Second * 70,
		Type:     entity.TargetTypeSlowQueryLog,
		URL:      "http://localhost:19000",
		Path:     constant.DefaultSlowQueryLogPath,
	},
}

func (r *targetRepository) SelectAll(_ context.Context) (entity.Targets, error) {
	return targets, nil
}

func (r *targetRepository) SelectByTargetType(_ context.Context, targetType entity.TargetType) (entity.Targets, error) {
	ret := make(entity.Targets, 0, len(targets))
	for _, target := range targets {
		if target.Type == targetType {
			ret = append(ret, target)
		}
	}
	return ret, nil
}

func (r *targetRepository) Update(_ context.Context, new entity.Targets) error {
	targets = new
	return nil
}
