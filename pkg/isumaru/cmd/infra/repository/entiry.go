package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/entity"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/constant"
	entity2 "github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/repository"
)

type entryRepository struct{}

func NewEntryRepository() repository.EntryRepository {
	return &entryRepository{}
}

func (r *entryRepository) SelectByEntryType(_ context.Context, entryType entity.EntryType) (entity2.Entries, error) {
	var dir string
	switch entryType {
	case entity.EntryTypeMysql:
		dir = constant.IsumaruSlowQueryLogDir
	default:
		return nil, errors.New("invalid entry type")
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	sort.Slice(dirEntries, func(i, j int) bool {
		return dirEntries[i].Name() > dirEntries[j].Name()
	})
	entries := make(entity2.Entries, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}
		unixStr := dirEntry.Name()
		unix, err := strconv.ParseInt(unixStr, 10, 64)
		if err != nil {
			return nil, err
		}

		targetDir := fmt.Sprintf("%s/%s", constant.IsumaruSlowQueryLogDir, unixStr)
		targetsDirEntries, err := os.ReadDir(targetDir)
		if err != nil {
			return nil, err
		}
		targetIDs := make([]string, 0, len(targetsDirEntries))
		for _, targetDirEntry := range targetsDirEntries {
			if targetDirEntry.IsDir() {
				continue
			}
			targetIDs = append(targetIDs, targetDirEntry.Name())
		}

		entries = append(entries, &entity2.Entry{
			ID:        unixStr,
			Type:      entryType,
			Time:      time.Unix(unix, 0),
			TargetIDs: targetIDs,
		})
	}

	return entries, nil
}
