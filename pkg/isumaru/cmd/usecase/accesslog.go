package usecase

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/service"
)

type AccessLogInteractor interface {
	// GetSlowRequests Fileからスローアクセスログを解析する
	GetSlowRequests(ctx context.Context, entryID, targetID string) (*SlowAccessInfo, error)
}

type accessLogInteractor struct {
	accessLogDirFormat string
	alpConfigPath      string
	collectService     service.CollectService
}

func NewAccessLogInteractor(
	accessLogDirFormat string,
	alpConfigPath string,
	collectService service.CollectService,
) AccessLogInteractor {
	return &accessLogInteractor{
		accessLogDirFormat,
		alpConfigPath,
		collectService,
	}
}

type SlowAccessInfo struct {
	Data      []byte
	TargetIDs []string
}

func (i *accessLogInteractor) GetSlowRequests(ctx context.Context, entryID, targetID string) (*SlowAccessInfo, error) {
	dir := fmt.Sprintf(i.accessLogDirFormat, entryID)
	path := fmt.Sprintf("%s/%s", dir, targetID)

	cmd := exec.Command("alp", "ltsv", "--config", i.alpConfigPath, "--format", "tsv", "--file", path)
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	targets, err := i.collectService.GetSucceededTargets(ctx, entryID, constant.TargetTypeAccessLog)
	if err != nil {
		return nil, err
	}

	targetIDs := make([]string, 0, len(targets))
	for _, target := range targets {
		targetIDs = append(targetIDs, target.ID)
	}

	return &SlowAccessInfo{
		Data:      data,
		TargetIDs: targetIDs,
	}, nil
}
