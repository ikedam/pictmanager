package handler

import (
	"context"
	"net/http"

	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/ikedam/pictmanager/pkg/service/tag"
	"github.com/labstack/echo/v4"
)

func Route(ctx context.Context, config *config.Config, g *echo.Group) error {
	c, err := tag.New(ctx, config)
	if err != nil {
		return err
	}
	g.GET("/", func(ec echo.Context) error {
		ctx := ec.Request().Context()
		tagList, err := c.GetTagList(ctx)
		if err != nil {
			return err
		}
		return ec.JSON(http.StatusOK, tagList)
	})
	return nil
}
