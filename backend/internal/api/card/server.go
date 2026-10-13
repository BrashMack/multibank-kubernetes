package card

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/masharpik/vtb-hack-back/pkg/reflector"
)

type Server struct {
	dependencies Dependencies
	conf         Config
}

func (s *Server) Init(ctx context.Context, cfg Config, api *echo.Group, deps Dependencies) error {
	s.dependencies = deps
	s.conf = cfg

	if err := reflector.CheckFieldsForNilInterfaces(s.dependencies); err != nil {
		return fmt.Errorf("reflector.CheckFieldsForNilInterfaces: %w", err)
	}

	if err := s.initRouting(api); err != nil {
		return fmt.Errorf("initRouting: %w", err)
	}

	return nil
}

func (s *Server) initRouting(api *echo.Group) error {
	api.GET("/cards", s.List)
	api.POST("/card/payment", s.Payment)
	api.GET("/card/:id", s.Get)

	return nil
}
