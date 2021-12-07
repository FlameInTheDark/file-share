package logs

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const LoggerKey = "logger"

// EchoLogger logger middleware to http requests
type EchoLogger struct {
	logger *zap.Logger
}

func NewEchoLogger(logger *zap.Logger) *EchoLogger {
	return &EchoLogger{logger: logger}
}

func (e *EchoLogger) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		scheme := "http"
		if c.Request().TLS != nil {
			scheme = "https"
		}
		e.logger.Debug(
			"HTTP request",
			zap.Time("timestamp", time.Now()),
			zap.String("http_method", c.Request().Method),
			zap.String("http_uri", c.Request().RequestURI),
			zap.String("http_addr", c.Request().RemoteAddr),
			zap.String("http_scheme", scheme),
			zap.String("http_agent", c.Request().UserAgent()),
		)
		//ctx := context.WithValue(c.Request().Context(), LoggerKey, e.logger)
		//c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}