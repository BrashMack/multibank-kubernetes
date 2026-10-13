package card

import (
	"context"

	"github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
	serviceModels "github.com/masharpik/vtb-hack-back/internal/service/card/models"
	paymentServiceModels "github.com/masharpik/vtb-hack-back/internal/service/payment/models"
)

type Dependencies struct {
	Service        Service
	JWTProvider    JWTProvider
	PaymentService PaymentService
}

// Service - интерфейс методов, необходимых серверу для запросов к сервисной логике по картам.
type Service interface {
	GetCardBalancesByLogin(ctx context.Context, req serviceModels.GetCardBalancesByLoginRequest) (
		*serviceModels.GetCardBalancesByLoginResponse, error)
	GetCardBalanceByLogin(ctx context.Context, req serviceModels.GetCardBalanceByLoginRequest) (
		*serviceModels.GetCardBalanceByLoginResponse, error)
}

type JWTProvider interface {
	ValidateToken(ctx context.Context, token string) (*models.Claims, error)
}

type PaymentService interface {
	Payment(ctx context.Context, req paymentServiceModels.PaymentRequest) (*paymentServiceModels.PaymentResponse, error)
}
