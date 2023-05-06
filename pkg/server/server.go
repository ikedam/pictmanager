package server

import (
	"context"

	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e      *echo.Echo
	config *config.Config
}

func New(ctx context.Context, config *config.Config) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	err := Route(ctx, config, e)
	if err != nil {
		return nil, err
	}
	return &Server{
		e:      e,
		config: config,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	return s.e.Start(":8080")
}
