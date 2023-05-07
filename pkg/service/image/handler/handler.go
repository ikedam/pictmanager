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
	g.GET("/:id", func(ec echo.Context) error {
		var request struct {
			ID string `param:"id"`
		}
		err := ec.Bind(&request)
		if err != nil {
			return rfc7807.BadRequest().WithError(err)
		}
		ctx := ec.Request().Context()
		image, err := c.GetImage(ctx, request.ID)
		if err != nil {
			return err
		}
		if image == nil {
			return rfc7807.NotFound()
		}
		return ec.JSON(http.StatusOK, image)
	})
	return nil
}

func AdminRoute(ctx context.Context, config *config.Config, g *echo.Group) error {
	c, err := image.New(ctx, config)
	if err != nil {
		return err
	}
	g.PUT("/:id", func(ec echo.Context) error {
		var id string
		err := echo.PathParamsBinder(ec).String("id", &id).BindError()
		if err != nil {
			return rfc7807.BadRequest().WithError(err)
		}
		ctx := ec.Request().Context()
		image, err := c.GetImage(ctx, id)
		if err != nil {
			return err
		}
		if image == nil {
			return rfc7807.NotFound()
		}
		priorTagList := image.TagList

		newImage := *image
		err = ec.Bind(&newImage)
		if err != nil {
			return rfc7807.BadRequest().WithError(err)
		}
		newImage.ID = image.ID
		newImage.CreateTime = image.CreateTime

		err = c.PutImageWithUpdatingTag(ctx, &newImage, priorTagList)
		if err != nil {
			return err
		}

		return ec.JSON(http.StatusOK, newImage)
	})
	return nil
}
