package handler

import (
	"context"
	"net/http"

	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/labstack/echo/v4"
)

func AdminRoute(ctx context.Context, config *config.Config, g *echo.Group) error {
	g.GET("/", func(ec echo.Context) error {
		return ec.JSON(http.StatusOK, true)
	})
	return nil
}
