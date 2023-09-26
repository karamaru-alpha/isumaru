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

	"github.com/karamaru-alpha/isumaru/web"
)

type Config struct {
	Port                  string
	SlowQueryLogDirFormat string
	AccessLogDirFormat    string
	SlpConfigPath         string
	AlpConfigPath         string
	AgentURL              string
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

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: web.BuildHTTPFS(),
		HTML5:      true,
	}))

	agentPort := port.NewAgentPort()
	entryRepository := repository.NewEntryRepository()
	targetRepository := repository.NewTargetRepository()

	collectService := service.NewCollectService(c.AccessLogDirFormat, c.SlowQueryLogDirFormat, agentPort, entryRepository, targetRepository)

	settingInteractor := usecase.NewSettingInteractor(c.AlpConfigPath, c.SlpConfigPath, targetRepository)
	collectInteractor := usecase.NewCollectInteractor(collectService, targetRepository, entryRepository)
	slowQueryLogInteractor := usecase.NewSlowQueryLogInteractor(c.SlpConfigPath, c.SlowQueryLogDirFormat, collectService)
	accessLogInteractor := usecase.NewAccessLogInteractor(c.AccessLogDirFormat, c.AlpConfigPath, collectService)

	settingHandler := handler.NewSettingHandler(settingInteractor)
	collectHandler := handler.NewCollectHandler(collectInteractor)
	slowQueryLogHandler := handler.NewSlowQueryLogHandler(slowQueryLogInteractor)
	accessLogHandler := handler.NewAccessLogHandler(accessLogInteractor)

	api := e.Group("/api")
	api.GET("/collect", collectHandler.Top)
	api.POST("/collect", collectHandler.Collect)
	api.GET("/setting", settingHandler.Top)
	api.POST("/setting/target", settingHandler.UpdateTargets)
	api.POST("/setting/slp", settingHandler.UpdateSlpConfig)
	api.POST("/setting/alp", settingHandler.UpdateAlpConfig)
	api.GET("/slowquerylog/:entryID/:targetID", slowQueryLogHandler.GetSlowQueries)
	api.GET("/accesslog/:entryID/:targetID", accessLogHandler.GetSlowRequests)

	if err := e.StartH2CServer(fmt.Sprintf(":%s", c.Port), &http2.Server{}); err != nil {
		slog.Error("failed to start web-agent server. err=%+v", err)
	}
}

type customValidator struct{}

func (v *customValidator) Validate(i interface{}) error {
	return validator.New().Struct(i)
}
