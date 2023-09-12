package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/port"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type MysqlInteractor interface {
	// Collect 競技サーバーに問い合わせ、N秒間の間に出力されたスロークエリログを集計する
	Collect(ctx context.Context, seconds int32, path string) error
}

type mysqlInteractor struct {
	port       port.AgentPort
	repository repository.SlowQueryLogRepository
}

func NewMysqlInteractor(
	port port.AgentPort,
	repository repository.SlowQueryLogRepository,
) MysqlInteractor {
	return &mysqlInteractor{
		port,
		repository,
	}
}

func (i *mysqlInteractor) Collect(ctx context.Context, seconds int32, path string) error {
	data, err := i.port.CollectSlowQueryLog(ctx, seconds, path)
	if err != nil {
		return err
	}

	now := time.Now()
	id := strconv.FormatInt(now.UnixNano(), 36)
	slowQueryLog := &entity.SlowQueryLog{
		ID:   id,
		Data: data,
	}

	if err := i.repository.Insert(ctx, slowQueryLog); err != nil {
		return err
	}

	return nil
}
