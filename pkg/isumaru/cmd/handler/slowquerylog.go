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
	interactor usecase.SlowQueryLogInteractor
}

func NewSlowQueryLogHandler(interactor usecase.SlowQueryLogInteractor) SlowQueryLogHandler {
	return &slowQueryLogHandler{interactor}
}

type GetSLowQueriesResponse struct {
	Data      []byte   `json:"data"`
	TargetIDs []string `json:"targetIDs"`
}

func (h *slowQueryLogHandler) GetSlowQueries(c echo.Context) error {
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
