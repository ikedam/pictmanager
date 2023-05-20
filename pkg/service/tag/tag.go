package tag

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/ikedam/pictmanager/pkg/model"
	"github.com/ikedam/pictmanager/pkg/simplestore"
	"github.com/pkg/errors"
)

func (c *Controller) GetTagList(ctx context.Context) ([]*model.Tag, error) {
	client, err := simplestore.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize firestore client")
	}
	var tagList []*model.Tag
	query := client.MustQuery(&tagList).
		Where("NormalizedTo", "==", "").
		OrderBy("Count", firestore.Desc)
	err = query.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query Image")
	}
	return tagList, nil
}
