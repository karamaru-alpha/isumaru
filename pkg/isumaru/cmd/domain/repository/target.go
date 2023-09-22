package repository

import (
	"context"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/entity"
)

type TargetRepository interface {
	SelectAll(ctx context.Context) (entity.Targets, error)
	SelectByTargetType(ctx context.Context, targetType entity.TargetType) (entity.Targets, error)
	Update(ctx context.Context, targets entity.Targets) error
}
