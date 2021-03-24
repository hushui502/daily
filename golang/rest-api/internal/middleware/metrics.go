package middleware

import (
	"github.com/labstack/echo/v4"
	"rest-api/pkg/metric"
	"time"
)

func (mw *MiddlewareManager) MetricsMiddleware(metrics metric.Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			var status int
			if err != nil {
				status = err.(*echo.HTTPError).Code
			} else {
				status = c.Response().Status
			}
			metrics.ObserveResponseTime(status, c.Request().Method, c.Path(), time.Since(start).Seconds())
			metrics.IncHits(status, c.Request().Method, c.Path())
			return err
		}
	}
}
