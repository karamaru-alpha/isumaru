package middleware

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/xcontext"
)

func NewContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			now := time.Now()
			ctx = xcontext.WithValue[xcontext.Now, time.Time](ctx, now)

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
