package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/usecase"
)

type MysqlHandler interface {
	GetSlowQueries(c echo.Context) error
}

type mysqlHandler struct {
	interactor usecase.MysqlInteractor
}

func NewMysqlHandler(interactor usecase.MysqlInteractor) MysqlHandler {
	return &mysqlHandler{interactor}
}

type GetSLowQueriesResponse struct {
	Data      []byte   `json:"data"`
	TargetIDs []string `json:"targetIDs"`
}

func (h *mysqlHandler) GetSlowQueries(c echo.Context) error {
	entryID := c.Param("entryID")
	targetID := c.Param("targetID")

	ctx := c.Request().Context()
	res, err := h.interactor.GetSlowQueries(ctx, entryID, targetID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &GetSLowQueriesResponse{
		Data:      res.Data,
		TargetIDs: res.TargetIDs,
	})
}
