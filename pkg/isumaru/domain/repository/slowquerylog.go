package repository

import (
	"context"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
)

type SlowQueryLogRepository interface {
	Insert(ctx context.Context, e *entity.SlowQueryLog) error
}
