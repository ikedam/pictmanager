package image

import (
	"context"

	"github.com/ikedam/pictmanager/pkg/config"
)

type Controller struct {
	config *config.Config
}

func New(ctx context.Context, config *config.Config) (*Controller, error) {
	return &Controller{
		config: config,
	}, nil
}
