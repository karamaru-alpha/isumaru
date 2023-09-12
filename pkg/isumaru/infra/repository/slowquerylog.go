package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/repository"
)

type slowQueryLogRepository struct{}

func NewSlowQueryLogRepository() repository.SlowQueryLogRepository {
	return &slowQueryLogRepository{}
}

func (r *slowQueryLogRepository) Insert(_ context.Context, e *entity.SlowQueryLog) error {
	if err := os.MkdirAll(constant.IsumaruSlowQueryLogDir, os.ModePerm); err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", constant.IsumaruSlowQueryLogDir, e.ID)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(e.Data); err != nil {
		return err
	}

	return nil
}
