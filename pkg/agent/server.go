package agent

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/felixge/fgprof"
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

func Serve(ctx context.Context, c *Config) {
	e := echo.New()
	e.HideBanner = true
	e.Validator = &customValidator{}

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	collectInteractor := usecase.NewCollectInteractor()
	collectHandler := handler.NewCollectHandler(collectInteractor)

	debug := e.Group("/debug")
	debug.GET("/pprof", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	debug.GET("/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	debug.GET("/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	debug.GET("/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	debug.GET("/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	debug.GET("/fgprof", echo.WrapHandler(fgprof.Handler()))
	debug.POST("/collect", collectHandler.Collect)

	go func() {
		if err := e.StartH2CServer(fmt.Sprintf(":%s", c.Port), &http2.Server{}); err != nil {
			slog.Error("failed to start web-agent server. err=%+v", err)
		}
	}()

	<-ctx.Done()
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
}

type customValidator struct{}

func (v *customValidator) Validate(i interface{}) error {
	return validator.New().Struct(i)
}
