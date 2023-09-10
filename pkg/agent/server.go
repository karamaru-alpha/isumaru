package agent

import (
	"fmt"

	"golang.org/x/exp/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"

	"github.com/karamaru-alpha/isumaru/pkg/agent/handler/mysql"
)

type Config struct {
	Port             string
	SlowQueryLogPath string
}

func Serve(c *Config) {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	mysqlHandler := mysql.New(c.SlowQueryLogPath)

	e.GET("/mysql/slowlog", mysqlHandler.TailSlowQueryLog)

	if err := e.StartH2CServer(fmt.Sprintf(":%s", c.Port), &http2.Server{}); err != nil {
		slog.Error("failed to start isumaru-agent server. err=%+v", err)
	}
}
