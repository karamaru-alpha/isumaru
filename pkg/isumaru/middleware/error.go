package middleware

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

func ErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			slog.Error(err.Error())
			return err
		}
		return nil
	}
}
