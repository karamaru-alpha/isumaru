package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/usecase"
)

type AccessLogHandler interface {
	GetSlowRequests(c echo.Context) error
}

type accessLogHandler struct {
	accessLogInteractor usecase.AccessLogInteractor
}

func NewAccessLogHandler(accessLogInteractor usecase.AccessLogInteractor) AccessLogHandler {
	return &accessLogHandler{accessLogInteractor}
}

type GetSlowRequestsResponse struct {
	Data      []byte   `json:"data"`
	TargetIDs []string `json:"targetIDs"`
}

func (h *accessLogHandler) GetSlowRequests(c echo.Context) error {
	entryID := c.Param("entryID")
	targetID := c.Param("targetID")

	ctx := c.Request().Context()
	res, err := h.accessLogInteractor.GetSlowRequests(ctx, entryID, targetID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &GetSlowRequestsResponse{
		Data:      res.Data,
		TargetIDs: res.TargetIDs,
	})
}
