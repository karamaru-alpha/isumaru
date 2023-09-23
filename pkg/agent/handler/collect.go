package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/agent/usecase"
)

type CollectHandler interface {
	Collect(c echo.Context) error
}

type collectHandler struct {
	interactor usecase.CollectInteractor
}

func NewCollectHandler(interactor usecase.CollectInteractor) CollectHandler {
	return &collectHandler{
		interactor,
	}
}

type CollectRequest struct {
	Seconds int32  `json:"seconds" validate:"min=1"`
	Path    string `json:"path" validate:"required"`
}

func (h *collectHandler) Collect(c echo.Context) error {
	r := &CollectRequest{}
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	ctx := c.Request().Context()
	reader, err := h.interactor.Collect(ctx, r.Seconds, r.Path)
	if err != nil {
		return err
	}

	return c.Stream(http.StatusOK, "application/json", reader)
}
