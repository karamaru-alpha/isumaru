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
	GetSlowQueries(ctx context.Context, entryID, targetID string) (*SlowQueryInfo, error)
}

type mysqlInteractor struct {
	collectService service.CollectService
}

func NewMysqlInteractor(collectService service.CollectService) MysqlInteractor {
	return &mysqlInteractor{collectService}
}

type SlowQueryInfo struct {
	Data      []byte
	TargetIDs []string
}

func (i *mysqlInteractor) GetSlowQueries(ctx context.Context, entryID, targetID string) (*SlowQueryInfo, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", constant.IsumaruEntryDir, entryID, constant.IsumaruSlowQueryLogDir, targetID)
	cmd := exec.CommandContext(ctx, "slp", "--config", constant.SlpConfigPath, "--format", "tsv", "--file", path)
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	targets, err := i.collectService.GetSucceededTargets(ctx, entryID, constant.TargetTypeSlowQueryLog)
	if err != nil {
		return nil, err
	}

	targetIDs := make([]string, 0, len(targets))
	for _, target := range targets {
		targetIDs = append(targetIDs, target.ID)
	}

	return &SlowQueryInfo{
		Data:      data,
		TargetIDs: targetIDs,
	}, nil
}
