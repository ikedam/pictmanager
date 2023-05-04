package model

import "time"

type Image struct {
	ID                 string     `json:"id"`
	TagList            []string   `json:"tagList"`
	Description        string     `json:"description"`
	MakingNote         string     `json:"makingNote"`
	ItemMask           int16      `json:"itemMask"`
	WithIkedam         bool       `json:"withIkedam" firestore:"-"`
	WithMinidam        bool       `json:"withMinidam" firestore:"-"`
	TwitterURL         string     `json:"twitterURL"`
	TweetComment       string     `json:"tweetComment"`
	PublishTime        time.Time  `json:"publishTime"`
	LastManualTagTime  *time.Time `json:"lastManualTagTime,omitempty"`
	LastMachineTagTime *time.Time `json:"lastMachineTagTime,omitempty"`
	CreateTime         time.Time  `json:"createTime"`
	UpdateTime         time.Time  `json:"updateTime"`
}
