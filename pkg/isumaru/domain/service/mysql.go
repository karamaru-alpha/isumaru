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

type MysqlService interface {
	Collect(ctx context.Context, entryID, targetID string) error
}

type mysqlService struct {
	agentPort        port.AgentPort
	entryRepository  repository.EntryRepository
	targetRepository repository.TargetRepository
}

func NewMysqlService(
	agentPort port.AgentPort,
	entryRepository repository.EntryRepository,
	targetRepository repository.TargetRepository,
) MysqlService {
	return &mysqlService{
		agentPort,
		entryRepository,
		targetRepository,
	}
}

func (s *mysqlService) Collect(ctx context.Context, entryID, targetID string) (err error) {
	// ターゲット情報を取得する
	target, err := s.targetRepository.LoadByPK(ctx, targetID, constant.TargetTypeSlowQueryLog)
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
			if err := s.entryRepository.UpdateTargetStatusFailure(ctx, entryID, targetID, constant.TargetTypeSlowQueryLog, err); err != nil {
				slog.Error(err.Error())
			}
			return
		}
		if err := s.entryRepository.UpdateTargetStatusSuccess(ctx, entryID, targetID, constant.TargetTypeSlowQueryLog); err != nil {
			slog.Error(err.Error())
		}
	}()

	// agentに問い合わせてスロークエリログのReaderを取得する
	reader, err := s.agentPort.CollectSlowQueryLog(ctx, target.URL, target.Path, target.Duration)
	if err != nil {
		return err
	}
	defer reader.Close()

	// スロークエリログをファイルに保存する
	dir := fmt.Sprintf("%s/%s/%s", constant.IsumaruEntryDir, entryID, constant.IsumaruSlowQueryLogDir)
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
