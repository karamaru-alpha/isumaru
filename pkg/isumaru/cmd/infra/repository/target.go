package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type targetRepository struct{}

func NewTargetRepository() repository.TargetRepository {
	return &targetRepository{}
}

var targets = entity.Targets{
	{
		ID:       "isu1",
		Duration: time.Second * 15,
		Type:     constant.TargetTypeSlowQueryLog,
		URL:      "http://localhost:19000",
		Path:     constant.DefaultSlowQueryLogPath,
	},
	{
		ID:       "isu1",
		Duration: time.Second * 15,
		Type:     constant.TargetTypeAccessLog,
		URL:      "http://localhost:19000",
		Path:     constant.DefaultAccessLogPath,
	},
}

func (r *targetRepository) LoadByPK(_ context.Context, targetID string, targetType constant.TargetType) (*entity.Target, error) {
	for _, target := range targets {
		if target.ID == targetID && target.Type == targetType {
			return target, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("target not found. id=%s, type=%d", targetID, targetType))
}

func (r *targetRepository) SelectAll(_ context.Context) (entity.Targets, error) {
	return targets, nil
}

func (r *targetRepository) Update(_ context.Context, new entity.Targets) error {
	targets = new
	return nil
}
