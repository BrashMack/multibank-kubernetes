package account

import (
	"context"

	"github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
	serviceModels "github.com/masharpik/vtb-hack-back/internal/service/card/models"
)

type Dependencies struct {
	Service     Service
	JWTProvider JWTProvider
}

// Service - интерфейс методов, необходимых серверу для запросов к сервисной логике по счетам.
type Service interface {
	GetAccountBalancesByLogin(ctx context.Context, req serviceModels.GetAccountBalancesByLoginRequest) (
		*serviceModels.GetAccountBalancesByLoginResponse, error)
	GetAccountBalanceByLogin(ctx context.Context, req serviceModels.GetAccountBalanceByLoginRequest) (
		*serviceModels.GetAccountBalanceByLoginResponse, error)
	GetAllTransactionsByLogin(ctx context.Context, req serviceModels.GetAllTransactionsByLoginRequest) (
		*serviceModels.GetAllTransactionsByLoginResponse, error)
}

type JWTProvider interface {
	ValidateToken(ctx context.Context, token string) (*models.Claims, error)
}
