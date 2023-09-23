package service

import (
	"context"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type EntryService interface {
	// GetEntries エントリを全件取得する
	GetEntries(ctx context.Context) (entity.Entries, error)
	// GetSucceededTargets 正常終了しているターゲットを取得する
	GetSucceededTargets(ctx context.Context, entityID string, targetType constant.TargetType) (entity.EntryTargets, error)
}

type entryService struct {
	entryRepository repository.EntryRepository
}

func NewEntryService(entryRepository repository.EntryRepository) EntryService {
	return &entryService{
		entryRepository,
	}
}

func (s *entryService) GetEntries(ctx context.Context) (entity.Entries, error) {
	entries, err := s.entryRepository.SelectAll(ctx)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *entryService) GetSucceededTargets(ctx context.Context, entityID string, targetType constant.TargetType) (entity.EntryTargets, error) {
	entry, err := s.entryRepository.LoadByPK(ctx, entityID)
	if err != nil {
		return nil, err
	}

	ret := make(entity.EntryTargets, 0, len(entry.Targets))
	for _, target := range entry.Targets {
		if target.StatusType != constant.EntryTargetStatusTypeSuccess || target.Type != targetType {
			continue
		}
		ret = append(ret, target)
	}

	return ret, nil
}
