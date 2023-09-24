package usecase

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/service"
)

type SlowQueryLogInteractor interface {
	// GetSlowQueries Fileからスロークエリログを解析する
	GetSlowQueries(ctx context.Context, entryID, targetID string) (*SlowQueryInfo, error)
}

type slowQueryLogInteractor struct {
	collectService service.CollectService
}

func NewSlowQueryLogInteractor(collectService service.CollectService) SlowQueryLogInteractor {
	return &slowQueryLogInteractor{collectService}
}

type SlowQueryInfo struct {
	Data      []byte
	TargetIDs []string
}

func (i *slowQueryLogInteractor) GetSlowQueries(ctx context.Context, entryID, targetID string) (*SlowQueryInfo, error) {
	dir := fmt.Sprintf(constant.IsumaruSlowQueryLogDirFormat, entryID)
	path := fmt.Sprintf("%s/%s", dir, targetID)
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
