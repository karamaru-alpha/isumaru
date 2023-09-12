package usecase

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/net/context"
)

type MysqlInteractor interface {
	// Collect N秒間waitし、その内に出力されたスロークエリログを返却する
	Collect(ctx context.Context, seconds int32, path string) (io.Reader, error)
}

type mysqlInteractor struct{}

func NewMysqlInteractor() MysqlInteractor {
	return &mysqlInteractor{}
}

func (i *mysqlInteractor) Collect(_ context.Context, seconds int32, path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open slow query log file. err=%+v", err)
	}

	start, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to seek slow query log file at start. err=%+v", err)
	}

	time.Sleep(time.Second * time.Duration(seconds))

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat slow query log file. err=%+v", err)
	}

	reader := io.LimitReader(file, stat.Size()-start)
	return reader, nil
}
