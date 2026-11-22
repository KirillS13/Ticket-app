package middlewares

import (
	"Ticket-app/internal/metrics"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func HTTPMetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == "/metrics" {
				return next(c)
			}
			start := time.Now()
			err := next(c)
			duration := time.Since(start).Seconds()
			status := c.Response().Status
			if err != nil {
				var httpErr *echo.HTTPError
				if errors.As(err, &httpErr) {
					status = httpErr.Code
				} else {
					status = http.StatusInternalServerError
				}
			}
			metrics.HTTPRequestDuration.WithLabelValues(c.Request().Method, c.Path(), strconv.Itoa(status)).Observe(duration)
			return err
		}
	}
}
