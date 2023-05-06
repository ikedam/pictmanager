package model

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type ImageItemMask int16

const (
	ImageItemIkedam  = ImageItemMask(0x01)
	ImageItemMinidam = ImageItemMask(0x02)
)

type Image struct {
	ID                 string        `json:"id"`
	ImageURL           string        `json:"imageURL" firestore:"-"`
	ThumbnailURL       string        `json:"thumbnailURL" firestore:"-"`
	TagList            []string      `json:"tagList"`
	Description        string        `json:"description"`
	MakingNote         string        `json:"makingNote"`
	ItemMask           ImageItemMask `json:"itemMask"`
	WithIkedam         bool          `json:"withIkedam" firestore:"-"`
	WithMinidam        bool          `json:"withMinidam" firestore:"-"`
	TwitterURL         string        `json:"twitterURL"`
	TweetComment       string        `json:"tweetComment"`
	PublishTime        time.Time     `json:"publishTime"`
	LastManualTagTime  *time.Time    `json:"lastManualTagTime,omitempty"`
	LastMachineTagTime *time.Time    `json:"lastMachineTagTime,omitempty"`
	CreateTime         time.Time     `json:"createTime"`
	UpdateTime         time.Time     `json:"updateTime"`
}

func (image *Image) OnLoadFromFirestore() error {
	if image.TagList == nil {
		image.TagList = []string{}
	}
	image.WithIkedam = ((image.ItemMask & ImageItemIkedam) != 0)
	image.WithMinidam = ((image.ItemMask & ImageItemMinidam) != 0)
	return nil
}

func (image *Image) FillURL(baseURL string) error {
	imageURL, err := url.JoinPath(baseURL, image.ID)
	if err != nil {
		return errors.Wrap(err, "failed to build URL for image")
	}
	thumbnailURL, err := url.JoinPath(baseURL, "thumbnail", image.ID)
	if err != nil {
		return errors.Wrap(err, "failed to build URL for thumbnail")
	}
	image.ImageURL = imageURL
	image.ThumbnailURL = thumbnailURL
	return nil
}
