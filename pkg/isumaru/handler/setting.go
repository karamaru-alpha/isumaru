package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/usecase"
)

type SettingHandler interface {
	Get(c echo.Context) error
	Update(c echo.Context) error
}

type settingHandler struct {
	interactor usecase.SettingInteractor
}

func NewSettingHandler(interactor usecase.SettingInteractor) SettingHandler {
	return &settingHandler{interactor}
}

type SettingGetResponse struct {
	Seconds            int32  `json:"seconds"`
	MainServerAddress  string `json:"mainServerAddress"`
	AccessLogPath      string `json:"accessLogPath"`
	MysqlServerAddress string `json:"mysqlServerAddress"`
	SlowQueryLogPath   string `json:"slowQueryLogPath"`
}

func (h *settingHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()

	setting, err := h.interactor.Get(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &SettingGetResponse{
		Seconds:            setting.Seconds,
		MainServerAddress:  setting.MainServerAddress,
		AccessLogPath:      setting.AccessLogPath,
		MysqlServerAddress: setting.MysqlServerAddress,
		SlowQueryLogPath:   setting.SlowQueryLogPath,
	})
}

type SettingUpdateRequest struct {
	Seconds            int32  `json:"seconds" validate:"min=1"`
	MainServerAddress  string `json:"mainServerAddress" validate:"required"`
	AccessLogPath      string `json:"accessLogPath" validate:"required"`
	MysqlServerAddress string `json:"mysqlServerAddress" validate:"required"`
	SlowQueryLogPath   string `json:"slowQueryLogPath" validate:"required"`
}

func (h *settingHandler) Update(c echo.Context) error {
	r := &SettingUpdateRequest{}
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.interactor.Update(ctx, r.Seconds, r.MainServerAddress, r.AccessLogPath, r.MysqlServerAddress, r.SlowQueryLogPath); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
