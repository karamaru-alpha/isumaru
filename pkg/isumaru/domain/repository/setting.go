package repository

import (
	"context"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
)

type SettingRepository interface {
	Load(ctx context.Context) (*entity.Setting, error)
	Update(ctx context.Context, e *entity.Setting) error
}
