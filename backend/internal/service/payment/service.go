package payment

import (
	"fmt"

	"github.com/masharpik/vtb-hack-back/pkg/reflector"
)

type Service struct {
	dependencies Dependencies
	conf         Config
}

func (s *Service) Init(cfg Config, deps Dependencies) error {
	s.dependencies = deps
	s.conf = cfg

	if err := reflector.CheckFieldsForNilInterfaces(s.dependencies); err != nil {
		return fmt.Errorf("reflector.CheckFieldsForNilInterfaces: %w", err)
	}

	return nil
}
