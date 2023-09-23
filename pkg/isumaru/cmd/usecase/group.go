package usecase

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/slog"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/service"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/xcontext"
)

type GroupInteractor interface {
	Top(ctx context.Context) (entity.Entries, error)
	Collect(ctx context.Context) error
}

type groupInteractor struct {
	entryService     service.EntryService
	mysqlService     service.MysqlService
	targetRepository repository.TargetRepository
	entryRepository  repository.EntryRepository
}

func NewGroupInteractor(
	entryService service.EntryService,
	mysqlService service.MysqlService,
	targetRepository repository.TargetRepository,
	entryRepository repository.EntryRepository,
) GroupInteractor {
	return &groupInteractor{
		entryService,
		mysqlService,
		targetRepository,
		entryRepository,
	}
}

func (i *groupInteractor) Top(ctx context.Context) (entity.Entries, error) {
	entries, err := i.entryService.GetEntries(ctx)
	if err != nil {
		return nil, nil
	}
	return entries, nil
}

func (i *groupInteractor) Collect(ctx context.Context) error {
	// 現在のターゲット一覧を取得する
	targets, err := i.targetRepository.SelectAll(ctx)
	if err != nil {
		return nil
	}

	// エントリを作成する
	now := xcontext.Value[xcontext.Now, time.Time](ctx)
	entryID := fmt.Sprintf("%d", now.Unix())
	entry := &entity.Entry{
		ID: entryID,
	}
	if err := i.entryRepository.Insert(ctx, entry); err != nil {
		return err
	}

	for _, target := range targets {
		target := target
		go func() {
			switch target.Type {
			case constant.TargetTypeSlowQueryLog:
				if err := i.mysqlService.Collect(ctx, entryID, target.ID); err != nil {
					slog.Error(err.Error())
				}
			}
		}()
	}
	return nil
}
