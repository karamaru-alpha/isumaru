package repository

import (
	"context"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
)

type EntryRepository interface {
	SelectByEntryType(ctx context.Context, entryType entity.EntryType) (entity.Entries, error)
}
