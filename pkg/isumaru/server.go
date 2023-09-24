package isumaru

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
	"golang.org/x/net/http2"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/handler"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/infra/port"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/infra/repository"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/cmd/usecase"
	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/service"
	xmiddleware "github.com/karamaru-alpha/isumaru/pkg/isumaru/middleware"
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
	e.Use(middleware.Gzip())
	e.Use(xmiddleware.ContextMiddleware)

	agentPort := port.NewAgentPort()
	entryRepository := repository.NewEntryRepository()
	targetRepository := repository.NewTargetRepository()

	collectService := service.NewCollectService(agentPort, entryRepository, targetRepository)

	settingInteractor := usecase.NewSettingInteractor(targetRepository)
	collectInteractor := usecase.NewCollectInteractor(collectService, targetRepository, entryRepository)
	slowQueryLogInteractor := usecase.NewSlowQueryLogInteractor(collectService)
	accessLogInteractor := usecase.NewAccessLogInteractor(collectService)

	settingHandler := handler.NewSettingHandler(settingInteractor)
	collectHandler := handler.NewCollectHandler(collectInteractor)
	slowQueryLogHandler := handler.NewSlowQueryLogHandler(slowQueryLogInteractor)
	accessLogHandler := handler.NewAccessLogHandler(accessLogInteractor)

	e.GET("/collect", collectHandler.Top)
	e.POST("/collect", collectHandler.Collect)
	e.GET("/setting", settingHandler.Top)
	e.POST("/setting/target", settingHandler.UpdateTargets)
	e.POST("/setting/slp", settingHandler.UpdateSlpConfig)
	e.POST("/setting/alp", settingHandler.UpdateAlpConfig)
	e.GET("/slowquerylog/:entryID/:targetID", slowQueryLogHandler.GetSlowQueries)
	e.GET("/accesslog/:entryID/:targetID", accessLogHandler.GetSlowRequests)

	if err := e.StartH2CServer(fmt.Sprintf(":%s", c.Port), &http2.Server{}); err != nil {
		slog.Error("failed to start web-agent server. err=%+v", err)
	}
}

type customValidator struct{}

func (v *customValidator) Validate(i interface{}) error {
	return validator.New().Struct(i)
}
