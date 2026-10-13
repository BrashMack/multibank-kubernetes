package payment

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	cardProviderModels "github.com/masharpik/vtb-hack-back/internal/service/card/models"
	consentProviderModels "github.com/masharpik/vtb-hack-back/internal/service/consent/models"
	"github.com/masharpik/vtb-hack-back/internal/service/payment/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

func (s *Service) Payment(ctx context.Context, req models.PaymentRequest) (*models.PaymentResponse, error) {
	if req.AccountIDFrom == req.AccountIDTo {
		return nil, models.ErrCannotPaymentToSameAccount
	}

	paymentAmount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", err, models.ErrInvalidPaymentAmount)
	}

	cardBalance, err := s.dependencies.CardService.GetAccountBalancesByLogin(ctx,
		cardProviderModels.GetAccountBalancesByLoginRequest{
			Login:         req.Login,
			AccountNumber: []string{req.AccountIDFrom, req.AccountIDTo},
		})
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.CardService.GetAccountBalancesByLogin: %w", err)
	}

	var (
		cardBalanceFrom models.CardBalance
		cardBalanceTo   models.CardBalance
	)

	for _, currCardBalance := range cardBalance.AccountBalances {
		if cardBalanceFrom.CardID == "" {
			if currCardBalance.Number == req.AccountIDFrom {
				cardBalanceFrom = models.CardBalance{
					CardID:     currCardBalance.Number,
					SchemeName: currCardBalance.SchemeName,
					Balance:    currCardBalance.Balance,
					Currency:   currCardBalance.Currency,
					Bank:       currCardBalance.Bank,
				}

				continue
			}
		}

		if cardBalanceTo.CardID == "" {
			if currCardBalance.Number == req.AccountIDTo {
				cardBalanceTo = models.CardBalance{
					CardID:     currCardBalance.Number,
					SchemeName: currCardBalance.SchemeName,
					Balance:    currCardBalance.Balance,
					Currency:   currCardBalance.Currency,
					Bank:       currCardBalance.Bank,
				}

				continue
			}
		}

		if cardBalanceFrom.CardID != "" && cardBalanceTo.CardID != "" {
			break
		}
	}

	if cardBalanceFrom.CardID == "" {
		return nil, models.ErrAccountFromNotFound
	}

	balanceFrom, err := strconv.ParseFloat(cardBalanceFrom.Balance, 64)
	if err != nil {
		return nil, fmt.Errorf("strconv.ParseFloat: %w", err)
	}

	if balanceFrom < paymentAmount {
		return nil, models.ErrNotEnoughMoney
	}

	consentResp, err := s.dependencies.ConsentService.CreateWriteConsent(ctx,
		consentProviderModels.CreateWriteConsentRequest{
			BankFrom:      cardBalanceFrom.Bank.ToLowerTrimSpace(),
			Login:         req.Login,
			ConsentType:   enums.Write,
			DebtorAccount: req.AccountIDFrom,
			Amount:        req.Amount,
		})
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.ConsentService.CreateWriteConsent: %w", err)
	}

	if cardBalanceTo.CardID != "" {
		var result *openapiProviderModels.CreatePaymentResponse

		result, err = s.dependencies.OpenApi.CreatePayment(ctx, consentResp.ConsentID,
			cardBalanceFrom.Bank.ToLowerTrimSpace(), cardBalanceTo.Bank.ToLowerTrimSpace(),
			req.Amount, cardBalanceFrom.Currency, req.Login,
			cardBalanceFrom.CardID, cardBalanceTo.CardID,
			cardBalanceFrom.SchemeName, cardBalanceTo.SchemeName)
		if err != nil {
			return nil, fmt.Errorf("s.dependencies.OpenApi.CreatePayment: %w", err)
		}

		return &models.PaymentResponse{
			PaymentID: result.Data.PaymentID,
			Status:    result.Data.Status,
		}, nil
	}

	var paymentErr error

	for _, currBank := range s.conf.Banks {
		var result *openapiProviderModels.CreatePaymentResponse

		shemeName := cardBalanceFrom.SchemeName
		if cardBalanceTo.SchemeName != "" {
			shemeName = cardBalanceTo.SchemeName
		}

		result, err = s.dependencies.OpenApi.CreatePayment(ctx, consentResp.ConsentID,
			cardBalanceFrom.Bank.ToLowerTrimSpace(), currBank.ToLowerTrimSpace(),
			req.Amount, cardBalanceFrom.Currency, req.Login,
			cardBalanceFrom.CardID, req.AccountIDTo,
			cardBalanceFrom.SchemeName, shemeName)
		if err != nil {
			paymentErr = errors.Join(paymentErr,
				fmt.Errorf("s.dependencies.OpenApi.CreatePayment: bank_to=%s: %w", currBank, err))

			continue
		}

		return &models.PaymentResponse{
			PaymentID: result.Data.PaymentID,
			Status:    result.Data.Status,
		}, nil
	}

	return nil, paymentErr
}
