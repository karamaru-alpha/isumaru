package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/entity"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/constant"
	entity2 "github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/port"
	repository2 "github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/repository"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/xcontext"
)

type MysqlInteractor interface {
	// Collect 競技サーバーに問い合わせ、スロークエリログをFileに保存する
	Collect(ctx context.Context) error
	// GetEntries Fileからスロークエリログの一覧を取得する
	GetEntries(ctx context.Context) (entity2.Entries, error)
	// GetSlowQueries Fileからスロークエリログを解析する
	GetSlowQueries(ctx context.Context, id, targetID string) ([]byte, error)
	// GetTargetIDs エントリの対象一覧を取得する
	GetTargetIDs(ctx context.Context, id string) ([]string, error)
}

type mysqlInteractor struct {
	agentPort        port.AgentPort
	entryRepository  repository2.EntryRepository
	targetRepository repository2.TargetRepository
}

func NewMysqlInteractor(
	agentPort port.AgentPort,
	entryRepository repository2.EntryRepository,
	targetRepository repository2.TargetRepository,
) MysqlInteractor {
	return &mysqlInteractor{
		agentPort,
		entryRepository,
		targetRepository,
	}
}

type AgentCollectSlowQueryLogRequest struct {
	Seconds int32  `json:"seconds"`
	Path    string `json:"path"`
}

func (i *mysqlInteractor) Collect(ctx context.Context) error {
	targets, err := i.targetRepository.SelectByTargetType(ctx, entity2.TargetTypeSlowQueryLog)
	if err != nil {
		return err
	}

	now := xcontext.Value[xcontext.Now, time.Time](ctx)
	unixTime := now.Unix()

	eg, ctx := errgroup.WithContext(ctx)
	for _, target := range targets {
		target := target
		eg.Go(func() error {
			// agentに問い合わせてスロークエリログのReaderを取得する
			reader, err := i.agentPort.CollectSlowQueryLog(ctx, target.URL, target.Path, target.Duration)
			if err != nil {
				return err
			}
			defer reader.Close()

			// スロークエリログをファイルに保存する
			dir := fmt.Sprintf("%s/%d", constant.IsumaruSlowQueryLogDir, unixTime)
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
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (i *mysqlInteractor) GetEntries(ctx context.Context) (entity2.Entries, error) {
	entries, err := i.entryRepository.SelectByEntryType(ctx, entity.EntryTypeMysql)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (i *mysqlInteractor) GetSlowQueries(ctx context.Context, id, targetID string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/%s", constant.IsumaruSlowQueryLogDir, id, targetID)
	cmd := exec.CommandContext(ctx, "slp", "--output", "standard", "--format", "tsv", "--file", path)
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (i *mysqlInteractor) GetTargetIDs(_ context.Context, id string) ([]string, error) {
	dir := fmt.Sprintf("%s/%s", constant.IsumaruSlowQueryLogDir, id)
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	targetIDs := make([]string, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		targetIDs = append(targetIDs, dirEntry.Name())
	}

	return targetIDs, nil
}
