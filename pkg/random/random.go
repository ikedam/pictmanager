package random

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/ikedam/pictmanager/pkg/log"
	"github.com/ikedam/pictmanager/pkg/model"
	"github.com/ikedam/pictmanager/pkg/service/image"
	"github.com/ikedam/pictmanager/pkg/simplestore"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Random struct {
	config *config.Config
}

func New(ctx context.Context, config *config.Config) (*Random, error) {
	return &Random{
		config: config,
	}, nil
}

func (r *Random) Scan(ctx context.Context) error {
	client, err := simplestore.New(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create firestore client")
	}

	log.Info(ctx, "scanning Image...")

	var imageList []*model.Image
	query := client.MustQuery(&imageList).
		OrderBy("PublishTime", firestore.Desc)
	count := 0
	err = query.Iter(ctx, func(v any) error {
		imageInfo := v.(*model.Image)
		imageInfo.Random = image.GetRandomValue()
		err := client.Put(ctx, imageInfo)
		if err != nil {
			return errors.Wrapf(err, "failed to set random value for %v", imageInfo.ID)
		}
		count = count + 1
		if count%100 == 0 {
			log.Debugf(ctx, "processed %v", count)
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to query Image")
	}
	log.Info(ctx, "Complete", zap.Int("processed", count))
	return nil
}
