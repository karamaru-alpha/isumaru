package handler

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/usecase"
)

type MysqlHandler interface {
	Collect(c echo.Context) error
	GetEntries(c echo.Context) error
	GetSlowQueries(c echo.Context) error
	GetSlowQueryTargets(c echo.Context) error
}

type mysqlHandler struct {
	interactor usecase.MysqlInteractor
}

func NewMysqlHandler(interactor usecase.MysqlInteractor) MysqlHandler {
	return &mysqlHandler{interactor}
}

func (h *mysqlHandler) Collect(c echo.Context) error {
	ctx := c.Request().Context()

	if err := h.interactor.Collect(ctx); err != nil {
		return nil
	}

	return nil
}

type MysqlGetEntriesResponse struct {
	Entries []*MysqlEntry `json:"entries"`
}

type MysqlEntry struct {
	ID        string   `json:"id"`
	UnixTime  int64    `json:"unixTime"`
	TargetIDs []string `json:"targetIDs"`
}

func (h *mysqlHandler) GetEntries(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.interactor.GetEntries(ctx)
	if err != nil {
		return nil
	}

	return c.JSON(http.StatusOK, toMysqlGetEntriesResponse(res))
}

func toMysqlGetEntriesResponse(entries entity.Entries) *MysqlGetEntriesResponse {
	response := &MysqlGetEntriesResponse{
		Entries: make([]*MysqlEntry, 0, len(entries)),
	}
	for _, e := range entries {
		response.Entries = append(response.Entries, &MysqlEntry{
			ID:        e.ID,
			UnixTime:  e.Time.Unix(),
			TargetIDs: e.TargetIDs,
		})
	}
	return response
}

func (h *mysqlHandler) GetSlowQueries(c echo.Context) error {
	id := c.Param("id")
	targetID := c.Param("targetID")

	ctx := c.Request().Context()
	data, err := h.interactor.GetSlowQueries(ctx, id, targetID)
	if err != nil {
		return err
	}

	return c.Stream(http.StatusOK, "application/json", bytes.NewBuffer(data))
}

type GetSlowQueryTargetsResponse struct {
	TargetIDs []string `json:"targetIDs"`
}

func (h *mysqlHandler) GetSlowQueryTargets(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	targetIDs, err := h.interactor.GetTargetIDs(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &GetSlowQueryTargetsResponse{
		TargetIDs: targetIDs,
	})
}
