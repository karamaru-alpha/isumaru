package usecase

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/service"
)

type MysqlInteractor interface {
	// GetSlowQueries Fileからスロークエリログを解析する
	GetSlowQueries(ctx context.Context, entryID, targetID string) ([]byte, error)
	// GetSucceededTargetIDs エントリの対象一覧を取得する
	GetSucceededTargetIDs(ctx context.Context, entryID string) ([]string, error)
}

type mysqlInteractor struct {
	collectService service.CollectService
}

func NewMysqlInteractor(collectService service.CollectService) MysqlInteractor {
	return &mysqlInteractor{collectService}
}

func (i *mysqlInteractor) GetSlowQueries(ctx context.Context, entryID, targetID string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", constant.IsumaruEntryDir, entryID, constant.IsumaruSlowQueryLogDir, targetID)
	cmd := exec.CommandContext(ctx, "slp", "--config", constant.SlpConfigPath, "--format", "tsv", "--file", path)
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (i *mysqlInteractor) GetSucceededTargetIDs(ctx context.Context, entryID string) ([]string, error) {
	targets, err := i.collectService.GetSucceededTargets(ctx, entryID, constant.TargetTypeSlowQueryLog)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(targets))
	for _, target := range targets {
		ret = append(ret, target.ID)
	}

	return ret, nil
}
