package payment

import (
	"context"

	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	cardProviderModels "github.com/masharpik/vtb-hack-back/internal/service/card/models"
	consentServiceModels "github.com/masharpik/vtb-hack-back/internal/service/consent/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

type Dependencies struct {
	OpenApi        OpenApi
	ConsentService ConsentService
	CardService    CardService
}

type OpenApi interface {
	CreatePayment(ctx context.Context, consentID string, bankFrom enums.Bank, bankTo enums.Bank, amount, currency,
		login, accountIDFrom, accountIDTo, schemeNameFrom, schemeNameTo string) (
		*openapiProviderModels.CreatePaymentResponse, error)
}

type ConsentService interface {
	CreateWriteConsent(ctx context.Context, req consentServiceModels.CreateWriteConsentRequest) (
		*consentServiceModels.CreateWriteConsentResponse, error)
}

type CardService interface {
	GetAccountBalancesByLogin(ctx context.Context, req cardProviderModels.GetAccountBalancesByLoginRequest) (
		*cardProviderModels.GetAccountBalancesByLoginResponse, error)
}
