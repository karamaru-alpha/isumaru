package usecase

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/slog"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/service"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/xcontext"
)

type CollectInteractor interface {
	Top(ctx context.Context) (entity.Entries, error)
	Collect(ctx context.Context) error
}

type collectInteractor struct {
	collectService   service.CollectService
	targetRepository repository.TargetRepository
	entryRepository  repository.EntryRepository
}

func NewCollectInteractor(
	collectService service.CollectService,
	targetRepository repository.TargetRepository,
	entryRepository repository.EntryRepository,
) CollectInteractor {
	return &collectInteractor{
		collectService,
		targetRepository,
		entryRepository,
	}
}

func (i *collectInteractor) Top(ctx context.Context) (entity.Entries, error) {
	entries, err := i.collectService.GetEntries(ctx)
	if err != nil {
		return nil, nil
	}
	return entries, nil
}

func (i *collectInteractor) Collect(ctx context.Context) error {
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
		// goroutineに処理を逃す
		go func() {
			if err := i.collectService.Collect(ctx, entryID, target.ID, target.Type); err != nil {
				slog.Error(err.Error())
			}
		}()
	}
	return nil
}
