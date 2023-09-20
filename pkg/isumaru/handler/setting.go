package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/entity"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/usecase"
)

type SettingHandler interface {
	Top(c echo.Context) error
	UpdateTargets(c echo.Context) error
	UpdateSlpConfig(e echo.Context) error
}

type settingHandler struct {
	interactor usecase.SettingInteractor
}

func NewSettingHandler(interactor usecase.SettingInteractor) SettingHandler {
	return &settingHandler{interactor}
}

type SettingTopResponse struct {
	Targets   []*SettingTarget `json:"targets"`
	SlpConfig string           `json:"slpConfig"`
}

type SettingTarget struct {
	ID       string `json:"id"`
	Type     int32  `json:"type"`
	URL      string `json:"url"`
	Path     string `json:"path"`
	Duration int    `json:"duration"`
}

func (h *settingHandler) Top(c echo.Context) error {
	ctx := c.Request().Context()

	info, err := h.interactor.Top(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &SettingTopResponse{
		Targets:   toSettingTargets(info.Targets),
		SlpConfig: info.SlpConfig,
	})
}

func toSettingTargets(targets entity.Targets) []*SettingTarget {
	ret := make([]*SettingTarget, 0, len(targets))
	for _, target := range targets {
		ret = append(ret, &SettingTarget{
			ID:       target.ID,
			Type:     int32(target.Type),
			URL:      target.URL,
			Path:     target.Path,
			Duration: int(target.Duration.Seconds()),
		})
	}
	return ret
}

func toSettingTargetEntities(targets []*SettingTarget) entity.Targets {
	ret := make(entity.Targets, 0, len(targets))
	for _, target := range targets {
		ret = append(ret, &entity.Target{
			ID:       target.ID,
			Type:     entity.TargetType(target.Type),
			URL:      target.URL,
			Path:     target.Path,
			Duration: time.Second * time.Duration(target.Duration),
		})
	}
	return ret
}

type SettingUpdateRequest struct {
	Targets []*SettingTarget `json:"targets" validate:"required"`
}

func (h *settingHandler) UpdateTargets(c echo.Context) error {
	r := &SettingUpdateRequest{}
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.interactor.UpdateTargets(ctx, toSettingTargetEntities(r.Targets)); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

type SettingUpdateSlpConfigRequest struct {
	SlpConfig string `json:"slpConfig" validate:"required"`
}

func (h *settingHandler) UpdateSlpConfig(c echo.Context) error {
	r := &SettingUpdateSlpConfigRequest{}
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.interactor.UpdateSlpConfig(ctx, r.SlpConfig); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
