package image

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/ikedam/pictmanager/pkg/model"
	"github.com/ikedam/pictmanager/pkg/simplestore"
	"github.com/pkg/errors"
)

func (c *Controller) GetImageList(ctx context.Context) ([]*model.Image, error) {
	client, err := simplestore.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize firestore client")
	}
	var imageList []*model.Image
	err = client.MustQuery(&imageList).
		OrderBy("PublishTime", firestore.Desc).
		Limit(10).Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query Image")
	}
	for _, image := range imageList {
		image.FillURL(c.config.GCSBaseURL())
	}
	return imageList, nil
}
