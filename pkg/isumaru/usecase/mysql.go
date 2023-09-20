package usecase

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type MysqlInteractor interface {
	// Collect 競技サーバーに問い合わせ、スロークエリログをFileに保存する
	Collect(ctx context.Context) error
	// GetEntries Fileからスロークエリログの一覧を取得する
	GetEntries(ctx context.Context) (entity.Entries, error)
	// GetSlowQueries Fileからスロークエリログを解析する
	GetSlowQueries(ctx context.Context, id, targetID string) ([]byte, error)
	// GetTargetIDs エントリの対象一覧を取得する
	GetTargetIDs(ctx context.Context, id string) ([]string, error)
}

type mysqlInteractor struct {
	entryRepository  repository.EntryRepository
	targetRepository repository.TargetRepository
}

func NewMysqlInteractor(
	entryRepository repository.EntryRepository,
	targetRepository repository.TargetRepository,
) MysqlInteractor {
	return &mysqlInteractor{
		entryRepository,
		targetRepository,
	}
}

type AgentCollectSlowQueryLogRequest struct {
	Seconds int32  `json:"seconds"`
	Path    string `json:"path"`
}

func (i *mysqlInteractor) Collect(ctx context.Context) error {
	targets, err := i.targetRepository.SelectByTargetType(ctx, entity.TargetTypeSlowQueryLog)
	if err != nil {
		return err
	}

	now := time.Now()
	unixTime := now.Unix()

	eg, ctx := errgroup.WithContext(ctx)
	for _, target := range targets {
		target := target
		eg.Go(func() error {
			// Agentに問い合わせてスロークエリログを取得する
			requestBody, err := json.Marshal(&AgentCollectSlowQueryLogRequest{
				Seconds: int32(target.Duration.Seconds()),
				Path:    target.Path,
			})
			if err != nil {
				return err
			}
			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				fmt.Sprintf("%s/mysql/collect", target.URL),
				bytes.NewBuffer(requestBody))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "gzip")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return errors.New(fmt.Sprintf("fail to collect slow query log. err=%+v", resp.Body))
			}
			if resp.ContentLength == 0 {
				return errors.New("slow query log is empty")
			}

			var reader io.ReadCloser
			switch resp.Header.Get("Content-Encoding") {
			case "gzip":
				gzipReader, err := gzip.NewReader(resp.Body)
				if err != nil {
					return err
				}
				defer gzipReader.Close()
				reader = gzipReader
			default:
				reader = resp.Body
			}

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

func (i *mysqlInteractor) GetEntries(ctx context.Context) (entity.Entries, error) {
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
