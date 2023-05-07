package image

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/ikedam/pictmanager/pkg/model"
	"github.com/ikedam/pictmanager/pkg/simplestore"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) GetImageList(ctx context.Context, count int, after string) ([]*model.Image, error) {
	client, err := simplestore.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize firestore client")
	}
	var imageList []*model.Image
	query := client.MustQuery(&imageList).
		OrderBy("PublishTime", firestore.Desc).
		Limit(count)
	if after == "" {
		err = query.Do(ctx)
	} else {
		err = query.DoAfterByID(ctx, after)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to query Image")
	}
	for _, image := range imageList {
		err := image.FillURL(c.config.GCSBaseURL())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to fill URL for %v", image.ID)
		}
	}
	return imageList, nil
}

func (c *Controller) GetImage(ctx context.Context, id string) (*model.Image, error) {
	client, err := simplestore.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize firestore client")
	}
	image := &model.Image{
		ID: id,
	}
	err = client.Get(ctx, image)
	if status.Code(err) == codes.NotFound {
		return nil, nil
	}
	err = image.FillURL(c.config.GCSBaseURL())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fill URL for %v", image.ID)
	}
	return image, nil
}
