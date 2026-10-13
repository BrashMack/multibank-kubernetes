package deposit

import (
	"context"
	"time"

	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	consentServiceModels "github.com/masharpik/vtb-hack-back/internal/service/consent/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

type Dependencies struct {
	OpenApi        OpenApi
	ConsentService ConsentService
	RedisApproved  RedisApproved
	RedisAwaiting  RedisAwaiting
}

type OpenApi interface {
	CreateAgreement(ctx context.Context, bank enums.Bank, login, consentID, productID, srcAccountNumber string,
		amount float64, termMonths int) (*openapiProviderModels.CreateAgreementResponse, error)
	GetAgreement(ctx context.Context, bank enums.Bank, login, consentID, agreementID string) (
		*openapiProviderModels.GetAgreementResponse, error)
	GetAgreements(ctx context.Context, bank enums.Bank, login, consentID string) (
		*openapiProviderModels.GetAgreementsResponse, error)
	GetProducts(ctx context.Context, bank enums.Bank, productType enums.AgreementType, consentID string) (
		*openapiProviderModels.GetProductsResponse, error)
	CreateReadConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string) (
		*openapiProviderModels.CreateReadConsentResponse, error)
	GetCreateConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank,
		login, consentID string) (*openapiProviderModels.GetCreateConsentResponse, error)
	GetConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string, consentID string) (
		*openapiProviderModels.GetConsentResponse, error)
}

type ConsentService interface {
	CreateCreateConsent(ctx context.Context, req consentServiceModels.CreateCreateConsentRequest) (
		*consentServiceModels.CreateCreateConsentResponse, error)
	CheckConsent(ctx context.Context, req consentServiceModels.CheckConsentRequest) error
}

type RedisApproved interface {
	SetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank, consentID string,
		ttl time.Duration) error
	GetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) (string, error)
	Delete(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) error
}

type RedisAwaiting interface {
	SetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank, consentID string,
		ttl time.Duration) error
	GetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) (string, error)
	Delete(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) error
}
