package repository

import (
	"context"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type settingRepository struct{}

func NewSettingRepository() repository.SettingRepository {
	return &settingRepository{}
}

var setting *entity.Setting

func (r *settingRepository) Load(_ context.Context) (*entity.Setting, error) {
	if setting == nil {
		return &entity.Setting{
			Seconds:            constant.DefaultCollectSeconds,
			MainServerAddress:  constant.DefaultMainServerAddress,
			AccessLogPath:      constant.DefaultAccessLogPath,
			MysqlServerAddress: constant.DefaultMysqlServerAddress,
			SlowQueryLogPath:   constant.DefaultSlowQueryLogPath,
		}, nil
	}

	return setting, nil
}

func (r *settingRepository) Update(_ context.Context, e *entity.Setting) error {
	setting = e
	return nil
}
