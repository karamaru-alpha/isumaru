package isumaru

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
	"golang.org/x/net/http2"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/constant"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/handler"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/infra/repository"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/usecase"
)

type Config struct {
	Port     string
	AgentURL string
}

func Serve(c *Config) {
	e := echo.New()
	e.HideBanner = true
	e.Validator = &customValidator{}

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	entryRepository := repository.NewEntryRepository()
	targetRepository := repository.NewTargetRepository()
	mysqlInteractor := usecase.NewMysqlInteractor(entryRepository, targetRepository)
	settingInteractor := usecase.NewSettingInteractor(targetRepository)

	mysqlHandler := handler.NewMysqlHandler(mysqlInteractor)
	settingHandler := handler.NewSettingHandler(settingInteractor)

	e.GET("/setting", settingHandler.Top)
	e.POST("/setting/target", settingHandler.UpdateTargets)
	e.POST("/setting/slp", settingHandler.UpdateSlpConfig)
	e.POST("/mysql/collect", mysqlHandler.Collect)
	e.GET("/mysql", mysqlHandler.GetEntries)
	e.GET("/mysql/:id/:targetID", mysqlHandler.GetSlowQueries)
	e.GET("/mysql/:id", mysqlHandler.GetSlowQueryTargets)

	if err := os.MkdirAll(constant.IsumaruSlowQueryLogDir, os.ModePerm); err != nil {
		panic(err)
	}

	if err := e.StartH2CServer(fmt.Sprintf(":%s", c.Port), &http2.Server{}); err != nil {
		slog.Error("failed to start web-agent server. err=%+v", err)
	}
}

type customValidator struct{}

func (v *customValidator) Validate(i interface{}) error {
	return validator.New().Struct(i)
}
