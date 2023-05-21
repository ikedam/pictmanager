package image

import (
	"context"
	"math/rand"

	"cloud.google.com/go/firestore"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ikedam/pictmanager/pkg/model"
	"github.com/ikedam/pictmanager/pkg/rfc7807"
	"github.com/ikedam/pictmanager/pkg/simplestore"
	"github.com/ikedam/pictmanager/pkg/util"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const RandomRange = 100

func GetRandomValue() int {
	return rand.Intn(RandomRange) + 1
}

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

func (c *Controller) CreateImage(ctx context.Context, image *model.Image) (*model.Image, error) {
	client, err := simplestore.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize firestore client")
	}
	now := util.GetCurrentTime()
	image.CreateTime = now
	image.UpdateTime = now
	image.Random = GetRandomValue()
	err = client.Put(ctx, image)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to put Image %v", image.ID)
	}
	return image, nil
}

func (c *Controller) PutImageWithUpdatingTag(ctx context.Context, image *model.Image, priorTagList []string, preserveTagTime bool) error {
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

	if !preserveTagTime && (newTagSet.Cardinality() > 0 || removeTagSet.Cardinality() > 0) {
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

func (c *Controller) GetImageForTagging(ctx context.Context) (*model.Image, error) {
	image, err := c.getImageForTaggingImpl(ctx)
	if err != nil {
		return nil, err
	}
	if image == nil {
		return nil, rfc7807.NotFound()
	}
	err = image.FillURL(c.config.GCSBaseURL())
	if err != nil {
		return nil, err
	}
	return image, nil
}

const pickRetry = 10

func (c *Controller) getImageForTaggingImpl(ctx context.Context) (*model.Image, error) {
	for count := 0; count < pickRetry; count++ {
		image, err := c.getImageNotTaggedRandomly(ctx)
		if err != nil {
			return nil, err
		}
		if image != nil {
			return image, nil
		}
	}

	// Then, try with already tagged image
	for count := 0; count < pickRetry; count++ {
		image, err := c.getImageTaggedRandomly(ctx)
		if err != nil {
			return nil, err
		}
		if image != nil {
			return image, nil
		}
	}

	return nil, nil
}

func (c *Controller) getImageNotTaggedRandomly(ctx context.Context) (*model.Image, error) {
	client, err := simplestore.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize firestore client")
	}
	var imageList []*model.Image
	randomValue := GetRandomValue()
	query := client.MustQuery(&imageList).
		Where("Random", "==", randomValue).
		Where("LastManualTagTime", "==", nil)
	err = query.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query Image")
	}
	if len(imageList) <= 0 {
		return nil, nil
	}
	return pickOneImage(imageList), nil
}

func (c *Controller) getImageTaggedRandomly(ctx context.Context) (*model.Image, error) {
	client, err := simplestore.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize firestore client")
	}
	var imageList []*model.Image
	randomValue := GetRandomValue()
	query := client.MustQuery(&imageList).
		Where("Random", "==", randomValue).
		OrderBy("LastManualTagTime", firestore.Asc).
		Limit(10)
	err = query.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query Image")
	}
	if len(imageList) <= 0 {
		return nil, nil
	}
	return pickOneImage(imageList), nil
}

func pickOneImage(imageList []*model.Image) *model.Image {
	idx := rand.Intn(len(imageList))
	return imageList[idx]
}
