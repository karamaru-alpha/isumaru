package usecase

import (
	"context"
	"io"
	"os"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type SettingInteractor interface {
	// Top 全設定を返却する
	Top(ctx context.Context) (*SettingTopInfo, error)
	// UpdateTargets ターゲット情報を更新する
	UpdateTargets(ctx context.Context, targets entity.Targets) error
	// UpdateSlpConfig slpのconfigファイルを更新する
	UpdateSlpConfig(ctx context.Context, config string) error
	// UpdateAlpConfig alpのconfigファイルを更新する
	UpdateAlpConfig(ctx context.Context, config string) error
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
	AlpConfig string
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

	alpConfigFile, err := os.Open(constant.AlpConfigPath)
	if err != nil {
		return nil, err
	}
	defer alpConfigFile.Close()
	alpConfigData, err := io.ReadAll(alpConfigFile)
	if err != nil {
		return nil, err
	}

	return &SettingTopInfo{
		Targets:   targets,
		SlpConfig: string(slpConfigData),
		AlpConfig: string(alpConfigData),
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

func (i *settingInteractor) UpdateAlpConfig(_ context.Context, config string) error {
	file, err := os.Create(constant.AlpConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(config); err != nil {
		return err
	}
	return nil
}
