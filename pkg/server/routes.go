package server

import (
	"context"

	"github.com/ikedam/pictmanager/pkg/config"
	image "github.com/ikedam/pictmanager/pkg/service/image/handler"
	login "github.com/ikedam/pictmanager/pkg/service/login/handler"
	session "github.com/ikedam/pictmanager/pkg/service/session/handler"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func Route(ctx context.Context, config *config.Config, e *echo.Echo) error {
	api := e.Group("/api")
	err := image.Route(ctx, config, api.Group("/image"))
	if err != nil {
		return errors.Wrap(err, "error in building /image")
	}
	admin := api.Group("/admin")
	err = session.AdminRoute(ctx, config, admin.Group("/session"))
	if err != nil {
		return errors.Wrap(err, "error in building /admin/session")
	}
	err = login.AdminRoute(ctx, config, admin.Group("/login"))
	if err != nil {
		return errors.Wrap(err, "error in building /admin/login")
	}
	return nil
}
