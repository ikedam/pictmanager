package handler

import (
	"context"
	"net/http"

	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/ikedam/pictmanager/pkg/model"
	"github.com/ikedam/pictmanager/pkg/rfc7807"
	"github.com/ikedam/pictmanager/pkg/service/image"
	"github.com/labstack/echo/v4"
)

func Route(ctx context.Context, config *config.Config, g *echo.Group) error {
	c, err := image.New(ctx, config)
	if err != nil {
		return err
	}
	g.GET("/", func(ec echo.Context) error {
		request := struct {
			Count int    `query:"count"`
			After string `query:"after"`
		}{
			Count: 10,
		}
		err := ec.Bind(&request)
		if err != nil {
			return rfc7807.BadRequest().WithError(err)
		}
		if request.Count > 100 {
			return rfc7807.BadRequest().WithDetailf("too many counts: %v", request.Count)
		}
		ctx := ec.Request().Context()
		imageList, err := c.GetImageList(ctx, request.Count, request.After)
		if err != nil {
			return err
		}
		if imageList == nil {
			imageList = []*model.Image{}
		}
		return ec.JSON(http.StatusOK, imageList)
	})
	return nil
}
