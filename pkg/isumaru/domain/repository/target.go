package repository

import (
	"context"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
)

type TargetRepository interface {
	LoadByPK(ctx context.Context, id string, targetType constant.TargetType) (*entity.Target, error)
	SelectAll(ctx context.Context) (entity.Targets, error)
	Update(ctx context.Context, targets entity.Targets) error
}
