package user

import (
	"context"
	jwtModels "github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
)

type Dependencies struct {
	Provider    Provider
	JWTProvider JWTProvider
}

// Provider - интерфейс методов, необходимых сервисной логике для запросов по пользователям в БД.
type Provider interface {
	ValidateCredentials(ctx context.Context, login, hashedPassword string) error
}

type JWTProvider interface {
	GenerateToken(ctx context.Context, login string) (string, error)
	ValidateToken(ctx context.Context, token string) (*jwtModels.Claims, error)
}
