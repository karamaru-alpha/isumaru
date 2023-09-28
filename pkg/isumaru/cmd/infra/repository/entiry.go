package repository

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type entryRepository struct {
	mu       sync.RWMutex
	entryMap map[string]*entity.Entry
}

func NewEntryRepository() repository.EntryRepository {
	return &entryRepository{
		entryMap: map[string]*entity.Entry{
			"1695550273": {
				ID: "1695550273",
				Targets: entity.EntryTargets{
					{
						Target: &entity.Target{
							ID:   "isu1",
							Type: constant.TargetTypeSlowQueryLog,
						},
						StatusType: constant.EntryTargetStatusTypeSuccess,
						Error:      nil,
					},
					{
						Target: &entity.Target{
							ID:   "isu1",
							Type: constant.TargetTypeAccessLog,
						},
						StatusType: constant.EntryTargetStatusTypeSuccess,
						Error:      nil,
					},
				},
			},
		},
	}
}

func (r *entryRepository) Insert(_ context.Context, entry *entity.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.entryMap[entry.ID]; ok {
		return errors.New("entry already exists")
	}
	r.entryMap[entry.ID] = entry
	return nil
}

func (r *entryRepository) LoadByPK(_ context.Context, entryID string) (*entity.Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, ok := r.entryMap[entryID]
	if !ok {
		return nil, errors.New("entry not found")
	}
	return entry, nil
}

func (r *entryRepository) SelectAll(_ context.Context) (entity.Entries, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ret := make(entity.Entries, 0, len(r.entryMap))
	for _, entry := range r.entryMap {
		ret = append(ret, entry)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].ID > ret[j].ID
	})
	return ret, nil
}

func (r *entryRepository) AddTarget(_ context.Context, entryID string, target *entity.EntryTarget) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.entryMap[entryID]; !ok {
		return errors.New("entry not found")
	}
	r.entryMap[entryID].Targets = append(r.entryMap[entryID].Targets, target)
	return nil
}

func (r *entryRepository) UpdateTargetStatusSuccess(_ context.Context, entryID, targetID string, targetType constant.TargetType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.entryMap[entryID]; !ok {
		return errors.New("entry not found")
	}
	for _, target := range r.entryMap[entryID].Targets {
		if target.ID == targetID && target.Type == targetType {
			target.StatusType = constant.EntryTargetStatusTypeSuccess
			return nil
		}
	}
	return errors.New("target not found")
}

func (r *entryRepository) UpdateTargetStatusFailure(_ context.Context, entryID, targetID string, targetType constant.TargetType, err error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.entryMap[entryID]; !ok {
		return errors.New("entry not found")
	}

	for _, target := range r.entryMap[entryID].Targets {
		if target.ID == targetID && target.Type == targetType {
			target.StatusType = constant.EntryTargetStatusTypeFailure
			target.Error = err
			return nil
		}
	}
	return errors.New("target not found")
}
