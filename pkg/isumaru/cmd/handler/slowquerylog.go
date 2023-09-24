package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/usecase"
)

type SlowQueryLogHandler interface {
	GetSlowQueries(c echo.Context) error
}

type slowQueryLogHandler struct {
	slowQueryLogInteractor usecase.SlowQueryLogInteractor
}

func NewSlowQueryLogHandler(slowQueryLogInteractor usecase.SlowQueryLogInteractor) SlowQueryLogHandler {
	return &slowQueryLogHandler{slowQueryLogInteractor}
}

type GetSlowQueriesResponse struct {
	Data      []byte   `json:"data"`
	TargetIDs []string `json:"targetIDs"`
}

func (h *slowQueryLogHandler) GetSlowQueries(c echo.Context) error {
	entryID := c.Param("entryID")
	targetID := c.Param("targetID")

	ctx := c.Request().Context()
	res, err := h.slowQueryLogInteractor.GetSlowQueries(ctx, entryID, targetID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &GetSlowQueriesResponse{
		Data:      res.Data,
		TargetIDs: res.TargetIDs,
	})
}
