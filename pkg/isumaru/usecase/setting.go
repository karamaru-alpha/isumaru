package usecase

import (
	"context"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type SettingInteractor interface {
	// Get 設定を取得する
	Get(ctx context.Context) (*entity.Setting, error)
	// Update 設定を更新する
	Update(ctx context.Context, seconds int32, mainServerAddress, accessLogPath, mysqlServerAddress, slowQueryLogPath string) error
}

type settingInteractor struct {
	repository repository.SettingRepository
}

func NewSettingInteractor(repository repository.SettingRepository) SettingInteractor {
	return &settingInteractor{
		repository,
	}
}

func (i *settingInteractor) Get(ctx context.Context) (*entity.Setting, error) {
	setting, err := i.repository.Load(ctx)
	if err != nil {
		return nil, err
	}

	return setting, nil
}

func (i *settingInteractor) Update(ctx context.Context, seconds int32, mainServerAddress, accessLogPath, mysqlServerAddress, slowQueryLogPath string) error {
	setting, err := i.repository.Load(ctx)
	if err != nil {
		return err
	}

	setting.Seconds = seconds
	setting.MainServerAddress = mainServerAddress
	setting.AccessLogPath = accessLogPath
	setting.MysqlServerAddress = mysqlServerAddress
	setting.SlowQueryLogPath = slowQueryLogPath
	if err := i.repository.Update(ctx, setting); err != nil {
		return err
	}

	return nil
}
