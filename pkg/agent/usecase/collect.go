package usecase

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/net/context"
)

type CollectInteractor interface {
	// Collect N秒間の内に出力されたログを返却する
	Collect(ctx context.Context, seconds int32, path string) (io.Reader, error)
}

type collectInteractor struct{}

func NewCollectInteractor() CollectInteractor {
	return &collectInteractor{}
}

func (i *collectInteractor) Collect(_ context.Context, seconds int32, path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file. err=%+v", err)
	}

	start, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to seek log file at start. err=%+v", err)
	}

	time.Sleep(time.Second * time.Duration(seconds))

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat log file. err=%+v", err)
	}

	reader := io.LimitReader(file, stat.Size()-start)
	return reader, nil
}
