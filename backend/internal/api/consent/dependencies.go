package consent

import (
	"context"

	jwtProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
	"github.com/masharpik/vtb-hack-back/internal/service/consent/models"
)

type Dependencies struct {
	Service     Service
	JWTProvider JWTProvider
}

// Service - интерфейс методов, необходимых серверу для запросов к сервисной логике по согласиям.
type Service interface {
	CreateReadConsent(ctx context.Context, request models.CreateReadConsentRequest) (
		*models.CreateReadConsentResponse, error)
	ListAwaiting(ctx context.Context, req models.ListAwaitingRequest) (*models.ListAwaitingResponse, error)
	ListActive(ctx context.Context, req models.ListActiveRequest) (*models.ListActiveResponse, error)
}

type JWTProvider interface {
	ValidateToken(ctx context.Context, token string) (*jwtProviderModels.Claims, error)
}
