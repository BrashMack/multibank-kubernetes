package deposit

import (
	"context"

	jwtProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
	"github.com/masharpik/vtb-hack-back/internal/service/deposit/models"
)

type Dependencies struct {
	Service     Service
	JWTProvider JWTProvider
}

type JWTProvider interface {
	ValidateToken(ctx context.Context, token string) (*jwtProviderModels.Claims, error)
}

type Service interface {
	CreateAgreement(ctx context.Context, req models.CreateAgreementRequest) (*models.CreateAgreementResponse, error)
	GetProductsByLogin(ctx context.Context, login string) (*models.GetProductsResponse, error)
	GetAgreementsByLogin(ctx context.Context, login string) (*models.GetAgreementsResponse, error)
	GetPriorityProduct(ctx context.Context) (*models.GetPriorityProductResponse, error)
}
