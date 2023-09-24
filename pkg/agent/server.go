package agent

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
	"golang.org/x/net/http2"

	"github.com/karamaru-alpha/isumaru/pkg/agent/cmd/handler"
	"github.com/karamaru-alpha/isumaru/pkg/agent/cmd/usecase"
)

type Config struct {
	Port string
}

func Serve(c *Config) {
	e := echo.New()
	e.HideBanner = true
	e.Validator = &customValidator{}

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	collectInteractor := usecase.NewCollectInteractor()
	collectHandler := handler.NewCollectHandler(collectInteractor)
	e.POST("/collect", collectHandler.Collect)

	if err := e.StartH2CServer(fmt.Sprintf(":%s", c.Port), &http2.Server{}); err != nil {
		slog.Error("failed to start web-agent server. err=%+v", err)
	}
}

type customValidator struct{}

func (v *customValidator) Validate(i interface{}) error {
	return validator.New().Struct(i)
}
