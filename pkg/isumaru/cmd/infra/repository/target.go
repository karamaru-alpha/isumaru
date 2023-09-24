package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type targetRepository struct {
	mu      sync.RWMutex
	targets entity.Targets
}

func NewTargetRepository() repository.TargetRepository {
	return &targetRepository{}
}

func (r *targetRepository) LoadByPK(_ context.Context, targetID string, targetType constant.TargetType) (*entity.Target, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, target := range r.targets {
		if target.ID == targetID && target.Type == targetType {
			return target, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("target not found. id=%s, type=%d", targetID, targetType))
}

func (r *targetRepository) SelectAll(_ context.Context) (entity.Targets, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.targets, nil
}

func (r *targetRepository) Update(_ context.Context, new entity.Targets) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.targets = new
	return nil
}
