package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func AdminRoute(ctx context.Context, config *config.Config, g *echo.Group) error {
	g.GET("/", func(ec echo.Context) error {
		ctx := ec.Request().Context()
		var request struct {
			Back string `query:"back"`
		}
		err := ec.Bind(&request)
		if err != nil {
			log.Warn(ctx, "error parsing request", zap.Error(err))
		}
		if !strings.HasPrefix(request.Back, "/") {
			request.Back = "/"
		}
		return ec.Redirect(http.StatusTemporaryRedirect, request.Back)
	})
	return nil
}
