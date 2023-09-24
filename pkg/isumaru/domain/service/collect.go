package service

import (
	"context"
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/slog"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/port"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type CollectService interface {
	// Collect 競技サーバーからログを取得する
	Collect(ctx context.Context, entryID, targetID string, targetType constant.TargetType) error
	// GetEntries エントリを全件取得する
	GetEntries(ctx context.Context) (entity.Entries, error)
	// GetSucceededTargets 正常終了しているターゲットを取得する
	GetSucceededTargets(ctx context.Context, entityID string, targetType constant.TargetType) (entity.EntryTargets, error)
}

type collectService struct {
	agentPort        port.AgentPort
	entryRepository  repository.EntryRepository
	targetRepository repository.TargetRepository
}

func NewCollectService(
	agentPort port.AgentPort,
	entryRepository repository.EntryRepository,
	targetRepository repository.TargetRepository,
) CollectService {
	return &collectService{
		agentPort,
		entryRepository,
		targetRepository,
	}
}

func (s *collectService) Collect(ctx context.Context, entryID, targetID string, targetType constant.TargetType) (err error) {
	// ターゲット情報を取得する
	target, err := s.targetRepository.LoadByPK(ctx, targetID, targetType)
	if err != nil {
		return nil
	}

	// エントリにターゲットを追加する
	if err := s.entryRepository.AddTarget(ctx, entryID, &entity.EntryTarget{
		Target:     target,
		StatusType: constant.EntryTargetStatusTypeProgress,
	}); err != nil {
		return err
	}
	// エントリターゲットの完了/エラー通知
	defer func() {
		if err != nil {
			if err := s.entryRepository.UpdateTargetStatusFailure(ctx, entryID, targetID, targetType, err); err != nil {
				slog.Error(err.Error())
			}
			return
		}
		if err := s.entryRepository.UpdateTargetStatusSuccess(ctx, entryID, targetID, targetType); err != nil {
			slog.Error(err.Error())
		}
	}()

	// agentに問い合わせてスロークエリログのReaderを取得する
	reader, err := s.agentPort.CollectLog(ctx, target.URL, target.Path, target.Duration)
	if err != nil {
		return err
	}
	defer reader.Close()

	// スロークエリログをファイルに保存する
	var dir string
	switch targetType {
	case constant.TargetTypeSlowQueryLog:
		dir = fmt.Sprintf(constant.IsumaruSlowQueryLogDirFormat, entryID)
	case constant.TargetTypeAccessLog:
		dir = fmt.Sprintf(constant.IsumaruAccessLogDirFormat, entryID)
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}
	path := fmt.Sprintf("%s/%s", dir, target.ID)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, reader); err != nil {
		return err
	}

	return nil
}

func (s *collectService) GetEntries(ctx context.Context) (entity.Entries, error) {
	entries, err := s.entryRepository.SelectAll(ctx)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *collectService) GetSucceededTargets(ctx context.Context, entityID string, targetType constant.TargetType) (entity.EntryTargets, error) {
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
