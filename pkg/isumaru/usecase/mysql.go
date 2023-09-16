package usecase

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type MysqlInteractor interface {
	// Collect 競技サーバーに問い合わせ、N秒間の間に出力されたスロークエリログをFileに保存する
	Collect(ctx context.Context, seconds int32, path string) error
	// GetEntries Fileからスロークエリログの一覧を取得する
	GetEntries(ctx context.Context) (entity.Entries, error)
	// GetSlowQueries Fileからスロークエリログを解析する
	GetSlowQueries(ctx context.Context, id string) ([]byte, error)
}

type mysqlInteractor struct {
	agentURL        string
	entryRepository repository.EntryRepository
}

func NewMysqlInteractor(
	agentURL string,
	entryRepository repository.EntryRepository,
) MysqlInteractor {
	return &mysqlInteractor{
		agentURL,
		entryRepository,
	}
}

type AgentCollectSlowQueryLogRequest struct {
	Seconds int32  `json:"seconds"`
	Path    string `json:"path"`
}

func (i *mysqlInteractor) Collect(ctx context.Context, seconds int32, slowQueryLogPath string) error {
	// Agentに問い合わせてスロークエリログを取得する
	requestBody, err := json.Marshal(&AgentCollectSlowQueryLogRequest{
		Seconds: seconds,
		Path:    slowQueryLogPath,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/mysql/collect", i.agentURL),
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
	now := time.Now()
	path := fmt.Sprintf("%s/%d", constant.IsumaruSlowQueryLogDir, now.UnixNano())
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, reader); err != nil {
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

func (i *mysqlInteractor) GetSlowQueries(ctx context.Context, id string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s", constant.IsumaruSlowQueryLogDir, id)

	cmd := exec.CommandContext(ctx, "slp", "--output", "standard", "--format", "tsv", "--file", path)
	res, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return res, nil
}
