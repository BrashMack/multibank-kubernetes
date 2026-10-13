package user

import (
	"context"
	"fmt"

	"github.com/masharpik/vtb-hack-back/internal/service/user/models"
)

func (s *Service) Login(ctx context.Context, request models.LoginRequest) (models.LoginResponse, error) {
	if err := s.dependencies.Provider.ValidateCredentials(ctx, request.Login, request.Password); err != nil {
		return models.LoginResponse{}, fmt.Errorf("s.dependencies.Provider.ValidateCredentials: %w", err)
	}

	token, err := s.dependencies.JWTProvider.GenerateToken(ctx, request.Login)
	if err != nil {
		return models.LoginResponse{}, fmt.Errorf("s.dependencies.JWTProvider.GenerateToken: %w", err)
	}

	return models.LoginResponse{
		Token: token,
	}, nil
}
