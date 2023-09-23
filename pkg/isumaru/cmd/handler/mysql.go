package handler

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/usecase"
)

type MysqlHandler interface {
	GetSlowQueries(c echo.Context) error
	GetSlowQueryTargets(c echo.Context) error
}

type mysqlHandler struct {
	interactor usecase.MysqlInteractor
}

func NewMysqlHandler(interactor usecase.MysqlInteractor) MysqlHandler {
	return &mysqlHandler{interactor}
}

func (h *mysqlHandler) GetSlowQueries(c echo.Context) error {
	entryID := c.Param("entryID")
	targetID := c.Param("targetID")

	ctx := c.Request().Context()
	data, err := h.interactor.GetSlowQueries(ctx, entryID, targetID)
	if err != nil {
		return err
	}

	return c.Stream(http.StatusOK, "application/json", bytes.NewBuffer(data))
}

type GetSlowQueryTargetsResponse struct {
	TargetIDs []string `json:"targetIDs"`
}

func (h *mysqlHandler) GetSlowQueryTargets(c echo.Context) error {
	entryID := c.Param("entryID")

	ctx := c.Request().Context()
	targetIDs, err := h.interactor.GetSucceededTargetIDs(ctx, entryID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &GetSlowQueryTargetsResponse{
		TargetIDs: targetIDs,
	})
}
