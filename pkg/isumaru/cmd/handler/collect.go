package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/usecase"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
)

type CollectHandler interface {
	Top(c echo.Context) error
	Collect(c echo.Context) error
}

type collectHandler struct {
	collectInteractor usecase.CollectInteractor
}

func NewCollectHandler(collectInteractor usecase.CollectInteractor) CollectHandler {
	return &collectHandler{collectInteractor}
}

func (h *collectHandler) Top(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.collectInteractor.Top(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &GroupTopResponse{
		Entries: toGroupEntries(res),
	})
}

func (h *collectHandler) Collect(c echo.Context) error {
	ctx := c.Request().Context()
	if err := h.collectInteractor.Collect(ctx); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

type GroupEntry struct {
	ID      string              `json:"id"`
	Targets []*GroupEntryTarget `json:"targets"`
}

type GroupEntryTarget struct {
	ID           string                         `json:"id"`
	Type         constant.TargetType            `json:"type"`
	StatusType   constant.EntryTargetStatusType `json:"statusType"`
	ErrorMessage string                         `json:"errorMessage"`
}

type GroupTopResponse struct {
	Entries []*GroupEntry `json:"entries"`
}

type GroupCollectResponse struct {
	Entry *GroupEntry `json:"entry"`
}

func toGroupEntry(e *entity.Entry) *GroupEntry {
	if e == nil {
		return nil
	}
	return &GroupEntry{
		ID:      e.ID,
		Targets: toGroupEntryTargets(e.Targets),
	}
}

func toGroupEntries(slice entity.Entries) []*GroupEntry {
	ret := make([]*GroupEntry, 0, len(slice))
	for _, entry := range slice {
		ret = append(ret, toGroupEntry(entry))
	}
	return ret
}

func toGroupEntryTargets(slice entity.EntryTargets) []*GroupEntryTarget {
	ret := make([]*GroupEntryTarget, 0, len(slice))
	for _, target := range slice {
		var errorMessage string
		if target.Error != nil {
			errorMessage = target.Error.Error()
		}
		ret = append(ret, &GroupEntryTarget{
			ID:           target.ID,
			Type:         target.Type,
			StatusType:   target.StatusType,
			ErrorMessage: errorMessage,
		})
	}
	return ret
}
