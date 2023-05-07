package image

import (
	"context"

	"cloud.google.com/go/firestore"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ikedam/pictmanager/pkg/model"
	"github.com/ikedam/pictmanager/pkg/simplestore"
	"github.com/ikedam/pictmanager/pkg/util"
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

func (c *Controller) GetImageListWithTag(ctx context.Context, tag string, count int, after string) ([]*model.Image, error) {
	client, err := simplestore.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize firestore client")
	}
	var imageList []*model.Image
	query := client.MustQuery(&imageList).
		Where("TagList", "array-contains", tag).
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

func (c *Controller) PutImageWithUpdatingTag(ctx context.Context, image *model.Image, priorTagList []string) error {
	client, err := simplestore.New(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to initialize firestore client")
	}
	now := util.GetCurrentTime()
	image.UpdateTime = now

	priorTagSet := mapset.NewSet[string]()
	priorTagSet.Append(priorTagList...)
	currentTagSet := mapset.NewSet[string]()
	currentTagSet.Append(image.TagList...)

	newTagSet := currentTagSet.Difference(priorTagSet)
	removeTagSet := priorTagSet.Difference(currentTagSet)

	if newTagSet.Cardinality() > 0 || removeTagSet.Cardinality() > 0 {
		image.LastManualTagTime = &now
	}

	err = client.Put(ctx, image)
	if err != nil {
		return err
	}
	err = image.FillURL(c.config.GCSBaseURL())
	if err != nil {
		return errors.Wrapf(err, "failed to fill URL for %v", image.ID)
	}
	for tag := range newTagSet.Iter() {
		tagEntry := &model.Tag{
			ID:    tag,
			Count: 0,
		}
		err := client.Get(ctx, tagEntry)
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return errors.Wrapf(err, "failed to get Tag %v", tag)
			}
		}
		tagEntry.Count += 1
		client.Put(ctx, tagEntry)
	}
	for tag := range removeTagSet.Iter() {
		tagEntry := &model.Tag{
			ID:    tag,
			Count: 0,
		}
		err := client.Get(ctx, tagEntry)
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return errors.Wrapf(err, "failed to get Tag %v", tag)
			}
			continue
		}
		tagEntry.Count -= 1
		client.Put(ctx, tagEntry)
	}
	return nil
}
