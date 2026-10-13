package card

import (
	"context"
	"errors"
	"fmt"
	"slices"

	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	redisApprovedModels "github.com/masharpik/vtb-hack-back/internal/provider/redis/approved/models"
	"github.com/masharpik/vtb-hack-back/internal/service/card/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

func (s *Service) GetAccounts(ctx context.Context, bank enums.Bank, login, consentID string) (
	*models.GetAccountsResponse, error) {
	accounts, getAccountsErr := s.dependencies.OpenApi.GetAccounts(ctx, bank, login, consentID)
	if getAccountsErr != nil {
		getAccountsErr = fmt.Errorf("s.dependencies.OpenApi.GetAccounts: %w", getAccountsErr)

		if errors.Is(getAccountsErr, openapiProviderModels.ErrConsentForbidden) {
			if err := s.dependencies.RedisApproved.Delete(ctx, enums.Read, login, bank); err != nil {
				return nil, fmt.Errorf("%w, s.dependencies.RedisApproved.Delete: consent_type=%s: %w",
					getAccountsErr, enums.Read, err)
			}

			if err := s.dependencies.RedisApproved.Delete(ctx, enums.Write, login, bank); err != nil {
				return nil, fmt.Errorf("%w, s.dependencies.RedisApproved.Delete: consent_type=%s: %w",
					getAccountsErr, enums.Write, err)
			}

			if err := s.dependencies.RedisApproved.Delete(ctx, enums.Create, login, bank); err != nil {
				return nil, fmt.Errorf("%w, s.dependencies.RedisApproved.Delete: consent_type=%s: %w",
					getAccountsErr, enums.Create, err)
			}
		}

		return nil, getAccountsErr
	}

	var accountBalances = &models.GetAccountsResponse{
		Accounts: make([]models.Account, 0, len(accounts.Data.Accounts)),
	}

	for _, currAccount := range accounts.Data.Accounts {
		var acc = models.Account{
			AccountID: currAccount.AccountID,
			Currency:  currAccount.Currency,
			Bank:      bank,
			Type:      currAccount.AccountSubType,
		}

		for _, currNumberCard := range currAccount.Cards {
			acc.Number = currNumberCard.Identification
			acc.SchemeName = currNumberCard.SchemeName

			break
		}

		accountBalances.Accounts = append(accountBalances.Accounts, acc)
	}

	return accountBalances, nil
}

func (s *Service) GetAccountBalances(ctx context.Context, bank enums.Bank, login, consentID string, accountID string,
	accountNumbers []string) (*models.GetAccountBalancesResponse, error) {
	accounts, err := s.GetAccounts(ctx, bank, login, consentID)
	if err != nil {
		return nil, fmt.Errorf("s.GetAccounts: %w", err)
	}

	var accountBalances = &models.GetAccountBalancesResponse{
		AccountBalances: make([]models.AccountBalance, 0, len(accounts.Accounts)),
	}

	for _, currAccount := range accounts.Accounts {
		if accountID != "" && accountID != currAccount.AccountID {
			continue
		}

		if accountNumbers != nil && len(accountNumbers) > 0 && !slices.Contains(accountNumbers, currAccount.Number) {
			continue
		}

		balances, getBalanceErr := s.dependencies.OpenApi.GetAccountBalance(ctx, bank, consentID, currAccount.AccountID)
		if getBalanceErr != nil {
			return nil, fmt.Errorf("s.dependencies.OpenApi.GetAccountBalance: account_id=%s: %w",
				currAccount.AccountID, getBalanceErr)
		}

		for _, balance := range balances.Data.Balance {
			if balance.Type == enums.InterimAvailableBalanceType {
				renamedBankName, ok := s.conf.BankToName[bank]
				if !ok {
					renamedBankName = bank
				}

				accountBalances.AccountBalances = append(accountBalances.AccountBalances, models.AccountBalance{
					AccountID:  currAccount.AccountID,
					Number:     currAccount.Number,
					Balance:    balance.Amount.Amount,
					Currency:   balance.Amount.Currency,
					SchemeName: currAccount.SchemeName,
					Bank:       renamedBankName,
					Type:       currAccount.Type,
				})

				break
			}
		}
	}

	return accountBalances, nil
}

func (s *Service) GetAccountBalancesByLogin(ctx context.Context, req models.GetAccountBalancesByLoginRequest) (
	*models.GetAccountBalancesByLoginResponse, error) {
	readConsents, err := s.getBankConsents(ctx, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.getBankConsents: %w", err)
	}

	var result = &models.GetAccountBalancesByLoginResponse{
		AccountBalances: make([]models.AccountBalance, 0),
	}

	// TODO: Здесь можно распалаллелить.
	for _, readConsent := range readConsents {
		resp, getBalancesErr := s.GetAccountBalances(ctx, readConsent.Bank, req.Login, readConsent.ConsentID,
			req.AccountID, req.AccountNumber)
		if getBalancesErr != nil {
			return nil, fmt.Errorf("s.GetAccountBalances: bank=%s, consent_id=%s: %w",
				readConsent.Bank, readConsent.ConsentID, getBalancesErr)
		}

		result.AccountBalances = append(result.AccountBalances, resp.AccountBalances...)
	}

	return result, nil
}

func (s *Service) GetAccountBalanceWithoutCards(ctx context.Context, bank enums.Bank, login, consentID,
	accountID string) (*models.GetAccountBalanceWithoutCardsResponse, error) {
	accounts, getBalancesErr := s.GetAccountBalances(ctx, bank, login, consentID,
		accountID, nil)
	if getBalancesErr != nil {
		return nil, fmt.Errorf("s.GetAccountBalances: bank=%s, consent_id=%s: %w",
			bank, consentID, getBalancesErr)
	}

	var result = new(models.GetAccountBalanceWithoutCardsResponse)

	for _, currAccount := range accounts.AccountBalances {
		result.AccountBalance = models.AccountBalance{
			AccountID: currAccount.AccountID,
			Number:    currAccount.Number,
			Balance:   currAccount.Balance,
			Currency:  currAccount.Currency,
			Bank:      currAccount.Bank,
			Type:      currAccount.Type,
		}

		break
	}

	if result.AccountBalance.AccountID == "" {
		return nil, models.ErrAccountNotFound
	}

	return result, nil
}

func (s *Service) GetAccountBalanceByLogin(ctx context.Context, req models.GetAccountBalanceByLoginRequest) (
	*models.GetAccountBalanceByLoginResponse, error) {
	readConsents, err := s.getBankConsents(ctx, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.getBankConsents: %w", err)
	}

	var result *models.GetAccountBalanceByLoginResponse

	// TODO: Здесь можно распалаллелить.
	for _, readConsent := range readConsents {
		resp, getBalancesErr := s.GetAccountBalanceWithoutCards(ctx, readConsent.Bank, req.Login, readConsent.ConsentID,
			req.AccountID)
		if getBalancesErr != nil {
			if errors.Is(getBalancesErr, models.ErrAccountNotFound) ||
				errors.Is(getBalancesErr, openapiProviderModels.ErrAccountNotFound) {
				continue
			}

			return nil, fmt.Errorf("s.GetAccountBalanceWithoutCards: bank=%s, consent_id=%s: %w",
				readConsent.Bank, readConsent.ConsentID, getBalancesErr)
		}

		result = &models.GetAccountBalanceByLoginResponse{
			AccountBalance: models.AccountBalanceData{
				AccountID: resp.AccountBalance.AccountID,
				Number:    resp.AccountBalance.Number,
				Balance:   resp.AccountBalance.Balance,
				Currency:  resp.AccountBalance.Currency,
				Bank:      resp.AccountBalance.Bank,
				Type:      resp.AccountBalance.Type,
			},
		}

		cards, getCardBalancesErr := s.GetCardBalances(ctx, resp.AccountBalance.Bank.ToLowerTrimSpace(),
			req.Login, readConsent.ConsentID)
		if getCardBalancesErr != nil {
			return nil, fmt.Errorf("s.GetCardBalances: bank=%s, consent_id=%s: %w",
				readConsent.Bank.ToLowerTrimSpace(), readConsent.ConsentID, getCardBalancesErr)
		}

		var cardsForAccount = make([]models.Card, 0)

		for _, currCard := range cards.CardBalances {
			if currCard.AccountNumber == result.AccountBalance.Number {
				cardsForAccount = append(cardsForAccount, models.Card{
					CardID:   currCard.CardID,
					Number:   currCard.CardNumber,
					Type:     currCard.CardType,
					CardName: currCard.CardName,
				})
			}
		}

		result.AccountBalance.NumberCards = cardsForAccount
	}

	return result, nil
}

func (s *Service) GetCardBalances(ctx context.Context, bank enums.Bank, login, consentID string) (
	*models.GetCardBalancesResponse, error) {
	cards, getCardsErr := s.dependencies.OpenApi.GetCards(ctx, bank, login, consentID)
	if getCardsErr != nil {
		getCardsErr = fmt.Errorf("s.dependencies.OpenApi.GetCards: %w", getCardsErr)

		if errors.Is(getCardsErr, openapiProviderModels.ErrConsentForbidden) {
			if err := s.dependencies.RedisApproved.Delete(ctx, enums.Read, login, bank); err != nil {
				return nil, fmt.Errorf("%w, s.dependencies.RedisApproved.Delete: consent_type=%s: %w",
					getCardsErr, enums.Read, err)
			}

			if err := s.dependencies.RedisApproved.Delete(ctx, enums.Write, login, bank); err != nil {
				return nil, fmt.Errorf("%w, s.dependencies.RedisApproved.Delete: consent_type=%s: %w",
					getCardsErr, enums.Write, err)
			}

			if err := s.dependencies.RedisApproved.Delete(ctx, enums.Create, login, bank); err != nil {
				return nil, fmt.Errorf("%w, s.dependencies.RedisApproved.Delete: consent_type=%s: %w",
					getCardsErr, enums.Create, err)
			}
		}

		return nil, getCardsErr
	}

	var cardBalances = &models.GetCardBalancesResponse{
		CardBalances: make([]models.CardInfoWithoutFullNumber, 0, len(cards.Data.Cards)),
	}

	var accountCurrency = make(map[string]struct {
		AccountID string
		Currency  string
	})

	for _, currCard := range cards.Data.Cards {
		accountCurrency[currCard.AccountNumber] = struct {
			AccountID string
			Currency  string
		}{
			AccountID: "",
			Currency:  "",
		}

		data := models.CardInfoWithoutFullNumber{
			CardID:         currCard.CardID,
			CardNumber:     currCard.CardNumber,
			CardType:       currCard.CardType,
			CardName:       currCard.CardName,
			HolderName:     currCard.HolderName,
			ExpiryMonth:    currCard.ExpiryMonth,
			ExpiryYear:     currCard.ExpiryYear,
			AccountNumber:  currCard.AccountNumber,
			AccountBalance: currCard.AccountBalance,
			Status:         currCard.Status,
			DailyLimit:     currCard.DailyLimit,
			MonthlyLimit:   currCard.MonthlyLimit,
			IssuedAt:       currCard.IssuedAt,
			Bank:           bank,
		}

		if renamedBankName, ok := s.conf.BankToName[data.Bank]; ok {
			data.Bank = renamedBankName
		}

		cardBalances.CardBalances = append(cardBalances.CardBalances, data)
	}

	accounts, err := s.GetAccounts(ctx, bank, login, consentID)
	if err != nil {
		return nil, fmt.Errorf("s.GetAccounts: %w", err)
	}

	for _, currAccount := range accounts.Accounts {
		if _, ok := accountCurrency[currAccount.Number]; ok && accountCurrency[currAccount.Number].AccountID == "" {
			accountCurrency[currAccount.Number] = struct {
				AccountID string
				Currency  string
			}{
				AccountID: currAccount.AccountID,
				Currency:  currAccount.Currency,
			}
		}
	}

	for i, currCardBalance := range cardBalances.CardBalances {
		cardBalances.CardBalances[i].AccountID = accountCurrency[currCardBalance.AccountNumber].AccountID
		cardBalances.CardBalances[i].AccountCurrency = accountCurrency[currCardBalance.AccountNumber].Currency
	}

	return cardBalances, nil
}

// GetCardBalancesByLogin - оставляем из общего списка счетов только счета карт и кредиток и формируем список карт.
func (s *Service) GetCardBalancesByLogin(ctx context.Context, req models.GetCardBalancesByLoginRequest) (
	*models.GetCardBalancesByLoginResponse, error) {
	readConsents, err := s.getBankConsents(ctx, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.getBankConsents: %w", err)
	}

	var result = &models.GetCardBalancesByLoginResponse{
		CardBalances: make([]models.CardInfoWithoutFullNumber, 0),
	}

	// TODO: Здесь можно распалаллелить.
	for _, readConsent := range readConsents {
		resp, getBalancesErr := s.GetCardBalances(ctx, readConsent.Bank, req.Login, readConsent.ConsentID)
		if getBalancesErr != nil {
			return nil, fmt.Errorf("s.GetCardBalances: bank=%s, consent_id=%s: %w",
				readConsent.Bank, readConsent.ConsentID, getBalancesErr)
		}

		result.CardBalances = append(result.CardBalances, resp.CardBalances...)
	}

	return result, nil
}

func (s *Service) GetCardBalance(ctx context.Context, bank enums.Bank, login, consentID, cardID string) (
	*models.GetCardBalanceResponse, error) {
	cardData, getCardsErr := s.dependencies.OpenApi.GetCard(ctx, bank, login, consentID, cardID)
	if getCardsErr != nil {
		getCardsErr = fmt.Errorf("s.dependencies.OpenApi.GetCards: %w", getCardsErr)

		if errors.Is(getCardsErr, openapiProviderModels.ErrConsentForbidden) {
			if err := s.dependencies.RedisApproved.Delete(ctx, enums.Read, login, bank); err != nil {
				return nil, fmt.Errorf("%w, s.dependencies.RedisApproved.Delete: consent_type=%s: %w",
					getCardsErr, enums.Read, err)
			}

			if err := s.dependencies.RedisApproved.Delete(ctx, enums.Write, login, bank); err != nil {
				return nil, fmt.Errorf("%w, s.dependencies.RedisApproved.Delete: consent_type=%s: %w",
					getCardsErr, enums.Write, err)
			}

			if err := s.dependencies.RedisApproved.Delete(ctx, enums.Create, login, bank); err != nil {
				return nil, fmt.Errorf("%w, s.dependencies.RedisApproved.Delete: consent_type=%s: %w",
					getCardsErr, enums.Create, err)
			}
		}

		return nil, getCardsErr
	}

	accounts, err := s.GetAccounts(ctx, bank, login, consentID)
	if err != nil {
		return nil, fmt.Errorf("s.GetAccounts: %w", err)
	}

	var (
		accountID, accountCurrency string
	)

	for _, currAccount := range accounts.Accounts {
		if currAccount.Number == cardData.CardData.AccountNumber {
			accountID = currAccount.AccountID
			accountCurrency = currAccount.Currency

			break
		}
	}

	var cardFullNumber = cardData.CardData.CardNumberFull
	if cardFullNumber == "" {
		cardFullNumber = cardData.CardData.CardNumber
	}

	var cardBalance = &models.GetCardBalanceResponse{
		CardBalance: models.CardInfo{
			CardID:          cardData.CardData.CardID,
			CardNumber:      cardData.CardData.CardNumber,
			CardFullNumber:  cardFullNumber,
			CardType:        cardData.CardData.CardType,
			CardName:        cardData.CardData.CardName,
			HolderName:      cardData.CardData.HolderName,
			ExpiryMonth:     cardData.CardData.ExpiryMonth,
			ExpiryYear:      cardData.CardData.ExpiryYear,
			AccountID:       accountID,
			AccountNumber:   cardData.CardData.AccountNumber,
			AccountBalance:  cardData.CardData.AccountBalance,
			AccountCurrency: accountCurrency,
			Status:          cardData.CardData.Status,
			DailyLimit:      cardData.CardData.DailyLimit,
			MonthlyLimit:    cardData.CardData.MonthlyLimit,
			IssuedAt:        cardData.CardData.IssuedAt,
			Bank:            bank,
		},
	}

	return cardBalance, nil
}

// GetCardBalanceByLogin - оставляем из общего списка счетов только счета карт и кредиток и формируем только
// нужную карту.
func (s *Service) GetCardBalanceByLogin(ctx context.Context, req models.GetCardBalanceByLoginRequest) (
	*models.GetCardBalanceByLoginResponse, error) {
	readConsents, err := s.getBankConsents(ctx, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.getBankConsents: %w", err)
	}

	var result = new(models.GetCardBalanceByLoginResponse)

	// TODO: Здесь можно распалаллелить.
	for _, readConsent := range readConsents {
		resp, getBalanceErr := s.GetCardBalance(ctx, readConsent.Bank, req.Login, readConsent.ConsentID, req.CardID)
		if getBalanceErr != nil {
			if errors.Is(getBalanceErr, openapiProviderModels.ErrCardNotFound) {
				continue
			}

			return nil, fmt.Errorf("s.GetCardBalance: bank=%s, consent_id=%s: %w",
				readConsent.Bank, readConsent.ConsentID, getBalanceErr)
		}

		if renamedBankName, ok := s.conf.BankToName[resp.CardBalance.Bank]; ok {
			resp.CardBalance.Bank = renamedBankName
		}

		result.CardBalance = resp.CardBalance
	}

	if result.CardBalance.CardID == "" {
		return nil, models.ErrCardNotFound
	}

	return result, nil
}

func (s *Service) getBankConsents(ctx context.Context, login string) ([]redisApprovedModels.ConsentData, error) {
	readConsents, err := s.dependencies.RedisApproved.GetBanksByLogin(ctx, enums.Read, login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisApproved.GetBanksByLogin: consent_type=%s: %w",
			enums.Read, err)
	}

	return readConsents, nil
}

func (s *Service) GetAllTransactionsByLogin(ctx context.Context, req models.GetAllTransactionsByLoginRequest) (
	*models.GetAllTransactionsByLoginResponse, error) {
	readConsents, err := s.dependencies.RedisApproved.GetBanksByLogin(ctx, enums.Read, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisApproved.GetBanksByLogin: %w", err)
	}

	result := &models.GetAllTransactionsByLoginResponse{
		Items: make([]models.Transaction, 0),
	}

	mapTransactions := func(txs []openapiProviderModels.Transaction) []models.Transaction {
		items := make([]models.Transaction, 0, len(txs))
		for _, t := range txs {
			info := t.TransactionInformation
			if info == "" {
				info = t.RemittanceInformation
			}
			if info == "" {
				info = t.Reference
			}

			var merchant *models.Merchant
			if t.Merchant != nil {
				merchant = &models.Merchant{
					MerchantID: t.Merchant.MerchantID,
					Name:       t.Merchant.Name,
					MccCode:    t.Merchant.MccCode,
					Category:   t.Merchant.Category,
					City:       t.Merchant.City,
					Country:    t.Merchant.Country,
					Address:    t.Merchant.Address,
				}
			}

			items = append(items, models.Transaction{
				TransactionID:          t.TransactionID,
				BookingDateTime:        t.BookingDateTime,
				ValueDateTime:          t.ValueDateTime,
				Amount:                 models.Amount{Amount: t.Amount.Amount, Currency: t.Amount.Currency},
				CreditDebitIndicator:   t.CreditDebitIndicator,
				Reference:              t.Reference,
				RemittanceInformation:  t.RemittanceInformation,
				TransactionInformation: info,
				CounterpartyAccount:    t.CounterpartyAccount,
				CounterpartyBank:       t.CounterpartyBank,
				Merchant:               merchant,
			})
		}
		return items
	}

	if req.AccountID != "" {
		for _, consent := range readConsents {
			accResp, err := s.dependencies.OpenApi.GetAccounts(ctx, consent.Bank, req.Login, consent.ConsentID)
			if err != nil {
				return nil, fmt.Errorf("s.dependencies.OpenApi.GetAccounts: %w", err)
			}

			for _, acc := range accResp.Data.Accounts {
				if acc.AccountID == req.AccountID {
					txResp, err := s.dependencies.OpenApi.GetAccountTransactions(
						ctx, consent.Bank, consent.ConsentID, req.AccountID, req.Login, req.Page, req.Limit,
					)
					if err != nil {
						if errors.Is(err, openapiProviderModels.ErrConsentForbidden) {
							return nil, fmt.Errorf("s.dependencies.OpenApi.GetAccountTransactions: %w", openapiProviderModels.ErrConsentForbidden)
						}
						return nil, fmt.Errorf("s.dependencies.OpenApi.GetAccountTransactions: %w", err)
					}

					result.Items = append(result.Items, mapTransactions(txResp.Data.Transaction)...)
					return result, nil
				}
			}
		}

		return result, nil
	}

	//for _, consent := range readConsents {
	//	accResp, err := s.dependencies.OpenApi.GetAccounts(ctx, consent.Bank, req.Login, consent.ConsentID)
	//	if err != nil {
	//		if errors.Is(err, openapiProviderModels.ErrConsentForbidden) {
	//			return nil, fmt.Errorf("bank=%s: %w", consent.Bank, openapiProviderModels.ErrConsentForbidden)
	//		}
	//		return nil, fmt.Errorf("s.dependencies.OpenApi.GetAccounts %w", err)
	//	}
	//
	//	for _, acc := range accResp.Data.Accounts {
	//		txResp, err := s.dependencies.OpenApi.GetAccountTransactions(
	//			ctx, consent.Bank, consent.ConsentID, acc.AccountID, req.Login, req.Page, req.Limit,
	//		)
	//		if err != nil {
	//			if errors.Is(err, openapiProviderModels.ErrConsentForbidden) {
	//				return nil, fmt.Errorf("s.dependencies.OpenApi.GetAccountTransactions: %w", openapiProviderModels.ErrConsentForbidden)
	//			}
	//			return nil, fmt.Errorf("s.dependencies.OpenApi.GetAccountTransactions: %w", err)
	//		}
	//		result.Items = append(result.Items, mapTransactions(txResp.Data.Transaction)...)
	//	}
	//}

	return result, nil
}
