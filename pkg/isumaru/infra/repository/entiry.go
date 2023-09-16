package repository

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type entryRepository struct{}

func NewEntryRepository() repository.EntryRepository {
	return &entryRepository{}
}

var entry *entity.Entry

func (r *entryRepository) SelectByEntryType(_ context.Context, entryType entity.EntryType) (entity.Entries, error) {
	var dirPath string
	switch entryType {
	case entity.EntryTypeMysql:
		dirPath = constant.IsumaruSlowQueryLogDir
	default:
		return nil, errors.New("invalid entry type")
	}

	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	entries := make(entity.Entries, 0, len(dirEntries))
	for _, e := range dirEntries {
		if e.IsDir() {
			continue
		}
		unixStr := e.Name()
		unix, err := strconv.ParseInt(unixStr, 10, 64)
		if err != nil {
			return nil, err
		}

		entries = append(entries, &entity.Entry{
			ID:   unixStr,
			Type: entryType,
			Time: time.Unix(unix, 0),
		})
	}

	return entries, nil
}
