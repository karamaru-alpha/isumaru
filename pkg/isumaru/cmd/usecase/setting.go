package usecase

import (
	"context"
	"io"
	"os"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/repository"
)

type SettingInteractor interface {
	// Top 全設定を返却する
	Top(ctx context.Context) (*SettingTopInfo, error)
	// UpdateTargets ターゲット情報を更新する
	UpdateTargets(ctx context.Context, targets entity.Targets) error
	// UpdateSlpConfig slpのconfigファイルを更新する
	UpdateSlpConfig(ctx context.Context, config string) error
}

type settingInteractor struct {
	targetRepository repository.TargetRepository
}

func NewSettingInteractor(targetRepository repository.TargetRepository) SettingInteractor {
	return &settingInteractor{
		targetRepository,
	}
}

type SettingTopInfo struct {
	Targets   entity.Targets
	SlpConfig string
}

func (i *settingInteractor) Top(ctx context.Context) (*SettingTopInfo, error) {
	targets, err := i.targetRepository.SelectAll(ctx)
	if err != nil {
		return nil, err
	}

	slpConfigFile, err := os.Open(constant.SlpConfigPath)
	if err != nil {
		return nil, err
	}
	defer slpConfigFile.Close()

	slpConfigData, err := io.ReadAll(slpConfigFile)
	if err != nil {
		return nil, err
	}

	return &SettingTopInfo{
		Targets:   targets,
		SlpConfig: string(slpConfigData),
	}, nil
}

func (i *settingInteractor) UpdateTargets(ctx context.Context, targets entity.Targets) error {
	if err := i.targetRepository.Update(ctx, targets); err != nil {
		return err
	}

	return nil
}

func (i *settingInteractor) UpdateSlpConfig(_ context.Context, config string) error {
	file, err := os.Create(constant.SlpConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(config); err != nil {
		return err
	}
	return nil
}
