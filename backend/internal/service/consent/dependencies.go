package consent

import (
	"context"
	"time"

	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	redisApprovedModels "github.com/masharpik/vtb-hack-back/internal/provider/redis/approved/models"
	redisAwaitingModels "github.com/masharpik/vtb-hack-back/internal/provider/redis/awaiting/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

type Dependencies struct {
	OpenApi       OpenApi
	RedisApproved RedisApproved
	RedisAwaiting RedisAwaiting
}

type OpenApi interface {
	CreateReadConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string) (
		*openapiProviderModels.CreateReadConsentResponse, error)
	CreateWriteConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login, amount,
		debtorAccount string) (*openapiProviderModels.CreateWriteConsentResponse, error)
	CreateCreateConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string) (
		*openapiProviderModels.CreateCreateConsentResponse, error)
	GetCreateConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank,
		login, consentID string) (*openapiProviderModels.GetCreateConsentResponse, error)
	GetConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string, consentID string) (
		*openapiProviderModels.GetConsentResponse, error)
}

type RedisApproved interface {
	SetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank, consentID string,
		ttl time.Duration) error
	Exists(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) (bool, error)
	GetBanksByLogin(ctx context.Context, typeConsent enums.ConsentType, login string) (
		[]redisApprovedModels.ConsentData, error)
	Delete(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) error
}

type RedisAwaiting interface {
	SetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank, requestID string,
		ttl time.Duration) error
	Exists(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) (bool, error)
	GetBanksByLogin(ctx context.Context, typeConsent enums.ConsentType, login string) (
		[]redisAwaitingModels.RequestData, error)
	Delete(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) error
}
