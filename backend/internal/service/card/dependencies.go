package card

import (
	"context"

	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	redisApprovedModels "github.com/masharpik/vtb-hack-back/internal/provider/redis/approved/models"
	consentServiceModels "github.com/masharpik/vtb-hack-back/internal/service/consent/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

type Dependencies struct {
	OpenApi        OpenApi
	ConsentService ConsentService
	RedisApproved  RedisApproved
}

type OpenApi interface {
	GetAccounts(ctx context.Context, bank enums.Bank, login, consentID string) (
		*openapiProviderModels.GetAccountsResponse, error)
	GetAccountBalance(ctx context.Context, bank enums.Bank, consentID, accountID string) (
		*openapiProviderModels.GetAccountBalanceResponse, error)
	GetCards(ctx context.Context, bank enums.Bank, login, consentID string) (
		*openapiProviderModels.GetCardsResponse, error)
	GetCard(ctx context.Context, bank enums.Bank, login, consentID, cardID string) (
		*openapiProviderModels.GetCardResponse, error)
	GetAccount(ctx context.Context, bank enums.Bank, login, consentID, accountID string) (
		*openapiProviderModels.GetAccountResponse, error)
	GetAccountTransactions(ctx context.Context, bank enums.Bank, consentID, accountID, login string, page, limit int) (
		*openapiProviderModels.GetAccountTransactionsResponse, error)
}

type ConsentService interface {
	ListConfirmed(ctx context.Context, req consentServiceModels.ListConfirmedRequest) (
		*consentServiceModels.ListConfirmedResponse, error)
	CreateWriteConsent(ctx context.Context, req consentServiceModels.CreateWriteConsentRequest) (
		*consentServiceModels.CreateWriteConsentResponse, error)
}

type RedisApproved interface {
	GetBanksByLogin(ctx context.Context, typeConsent enums.ConsentType, login string) (
		[]redisApprovedModels.ConsentData, error)
	Delete(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) error
}
