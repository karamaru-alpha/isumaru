package repository

import (
	"context"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
)

type EntryRepository interface {
	Insert(ctx context.Context, entry *entity.Entry) error
	LoadByPK(ctx context.Context, entryID string) (*entity.Entry, error)
	SelectAll(ctx context.Context) (entity.Entries, error)
	AddTarget(ctx context.Context, entryID string, target *entity.EntryTarget) error
	UpdateTargetStatusSuccess(ctx context.Context, entryID, targetID string, targetType constant.TargetType) error
	UpdateTargetStatusFailure(ctx context.Context, entryID, targetID string, targetType constant.TargetType, err error) error
}
