package echolog

import (
	"time"

	"github.com/ikedam/pictmanager/pkg/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Middleware() echo.MiddlewareFunc {
	requestLogger := log.NewLogger(zap.String("type", "request"))
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			ctx := log.CtxWithLogger(
				req.Context(),
				zap.String("type", "application"),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
			)
			req = req.WithContext(ctx)
			c.SetRequest(req)

			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			logFields := []zapcore.Field{
				zap.Int("status", res.Status),
				zap.String("method", req.Method),
				zap.String("remote_ip", c.RealIP()),
				zap.String("host", req.Host),
				zap.String("uri", req.RequestURI),
				zap.Float32("request_time", float32(stop.Sub(start).Microseconds())/1000000.0),
				zap.String("referer", req.Referer()),
				zap.String("user_agent", req.UserAgent()),
			}
			level := zap.InfoLevel
			if res.Status >= 500 {
				level = zap.ErrorLevel
			}
			if err != nil {
				logFields = append(logFields, zap.Error(err))
			}
			requestLogger.Log(
				level,
				req.Method+" "+req.RequestURI,
				logFields...,
			)
			return nil
		}
	}
}
