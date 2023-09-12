package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/usecase"
)

type MysqlHandler interface {
	Collect(c echo.Context) error
}

type mysqlHandler struct {
	interactor usecase.MysqlInteractor
}

func NewMysqlHandler(interactor usecase.MysqlInteractor) MysqlHandler {
	return &mysqlHandler{interactor}
}

func (h *mysqlHandler) Collect(c echo.Context) error {
	ctx := c.Request().Context()

	if err := h.interactor.Collect(ctx, 10, "testdata/slow-query.log"); err != nil {
		return nil
	}

	return nil
}
