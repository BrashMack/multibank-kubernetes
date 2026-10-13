package deposit

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	redisApprovedProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/redis/approved/models"
	redisAwaitingProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/redis/awaiting/models"
	consentServiceModels "github.com/masharpik/vtb-hack-back/internal/service/consent/models"
	"github.com/masharpik/vtb-hack-back/internal/service/deposit/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

func (s *Service) CreateAgreement(ctx context.Context, req models.CreateAgreementRequest) (
	*models.CreateAgreementResponse, error) {
	consentID, err := s.getCreateConsent(ctx, enums.Create, req.Bank, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.getCreateConsent: %w", err)
	}

	agreementAmount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, fmt.Errorf("strconv.ParseFloat: %w: %w", err, models.ErrInvalidAgreementAmount)
	}

	createAgreementResp, err := s.dependencies.OpenApi.CreateAgreement(ctx, req.Bank, req.Login, consentID,
		req.ProductID, req.SourceAccountID, agreementAmount, req.TermMonths)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.OpenApi.CreateAgreement: %w", err)
	}

	getAgreementResp, err := s.dependencies.OpenApi.GetAgreement(ctx, req.Bank, req.Login, consentID,
		createAgreementResp.Data.AgreementID)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.OpenApi.GetAgreement: %w", err)
	}

	renamedBank, ok := s.conf.BankToName[req.Bank]
	if !ok {
		renamedBank = req.Bank
	}

	return &models.CreateAgreementResponse{
		AgreementID:   createAgreementResp.Data.AgreementID,
		ProductID:     createAgreementResp.Data.ProductID,
		ProductName:   createAgreementResp.Data.ProductName,
		Amount:        fmt.Sprintf("%0.2f", getAgreementResp.Data.AccountBalance),
		AccountNumber: getAgreementResp.Data.AccountNumber,
		Bank:          renamedBank,
		Currency:      s.conf.DefaultCurrency,
		InterestRate:  fmt.Sprintf("%0.1f", getAgreementResp.Data.InterestRate),
	}, nil
}

func (s *Service) GetProducts(ctx context.Context, req models.GetProductsRequest) (*models.GetProductsResponse, error) {
	readConsent, err := s.getReadConsent(ctx, enums.Read, req.Bank, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.getReadConsent: %w", err)
	}

	resp, err := s.dependencies.OpenApi.GetProducts(ctx, req.Bank, enums.DepositAgreementType, readConsent)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.OpenApi.GetProducts: %w", err)
	}

	var products = make([]models.Product, 0, len(resp.Data.Product))

	for _, product := range resp.Data.Product {
		renamedBank, ok := s.conf.BankToName[req.Bank]
		if !ok {
			renamedBank = req.Bank
		}

		products = append(products, models.Product{
			ProductID:    product.ProductID,
			ProductName:  product.ProductName,
			Bank:         renamedBank,
			Currency:     s.conf.DefaultCurrency,
			InterestRate: product.InterestRate,
			MinAmount:    product.MinAmount,
			MaxAmount:    product.MaxAmount,
			TermMonths:   product.TermMonths,
		})
	}

	return &models.GetProductsResponse{
		Products: products,
	}, nil
}

func (s *Service) GetProductsByLogin(ctx context.Context, login string) (*models.GetProductsResponse, error) {
	var resultProducts = make([]models.Product, 0)

	for _, bank := range s.conf.Banks {
		products, err := s.GetProducts(ctx, models.GetProductsRequest{
			Bank:  bank,
			Login: login,
		})
		if err != nil {
			return nil, fmt.Errorf("s.GetProducts: %w", err)
		}

		resultProducts = append(resultProducts, products.Products...)
	}

	return &models.GetProductsResponse{
		Products: resultProducts,
	}, nil
}

func (s *Service) GetAgreements(ctx context.Context, req models.GetAgreementsRequest) (*models.GetAgreementsResponse,
	error) {
	createConsent, err := s.getCreateConsent(ctx, enums.Create, req.Bank, req.Login)
	switch {
	case errors.Is(err, consentServiceModels.ErrConsentNotActive):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("s.getCreateConsent: %w", err)
	}

	resp, err := s.dependencies.OpenApi.GetAgreements(ctx, req.Bank, req.Login, createConsent)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.OpenApi.GetAgreements: %w", err)
	}

	var agreements = make([]models.Agreement, 0, len(resp.Data))

	for _, agreement := range resp.Data {
		if agreement.Status != enums.ActiveAgreementStatus {
			continue
		}

		if agreement.ProductType != enums.DepositAgreementType {
			continue
		}

		agreementData, getAgreementErr := s.dependencies.OpenApi.GetAgreement(ctx, req.Bank, req.Login,
			createConsent, agreement.AgreementID)
		if getAgreementErr != nil {
			return nil, fmt.Errorf("s.dependencies.OpenApi.GetAgreement: %w", getAgreementErr)
		}

		renamedBank, ok := s.conf.BankToName[req.Bank]
		if !ok {
			renamedBank = req.Bank
		}

		agreements = append(agreements, models.Agreement{
			AgreementID:   agreement.AgreementID,
			ProductID:     agreement.ProductID,
			ProductName:   agreement.ProductName,
			Bank:          renamedBank,
			Currency:      s.conf.DefaultCurrency,
			InterestRate:  fmt.Sprintf("%0.1f", agreementData.Data.InterestRate),
			Amount:        fmt.Sprintf("%0.2f", agreementData.Data.AccountBalance),
			AccountNumber: agreement.AccountNumber,
		})
	}

	return &models.GetAgreementsResponse{
		Agreements: agreements,
	}, nil
}

func (s *Service) GetAgreementsByLogin(ctx context.Context, login string) (*models.GetAgreementsResponse, error) {
	var resultAgreement = make([]models.Agreement, 0)

	for _, bank := range s.conf.Banks {
		agreements, err := s.GetAgreements(ctx, models.GetAgreementsRequest{
			Bank:  bank,
			Login: login,
		})
		if err != nil {
			return nil, fmt.Errorf("s.GetAgreements: %w", err)
		}

		if agreements != nil {
			resultAgreement = append(resultAgreement, agreements.Agreements...)
		}
	}

	return &models.GetAgreementsResponse{
		Agreements: resultAgreement,
	}, nil
}

func (s *Service) GetPriorityProduct(ctx context.Context) (*models.GetPriorityProductResponse, error) {
	if s.conf.Priority.ProductID == "" {
		return &models.GetPriorityProductResponse{
			ProductID:    "deposit-reliable",
			ProductName:  "Вклад \"Надежный\"",
			Bank:         "SmartBank",
			Currency:     "RUB",
			InterestRate: "8.50",
			MinAmount:    "10000.00",
			MaxAmount:    "10000000.00",
			TermMonths:   12,
		}, nil
	}

	return &models.GetPriorityProductResponse{
		ProductID:    s.conf.Priority.ProductID,
		ProductName:  s.conf.Priority.ProductName,
		Bank:         s.conf.Priority.Bank,
		Currency:     s.conf.Priority.Currency,
		InterestRate: s.conf.Priority.InterestRate,
		MinAmount:    s.conf.Priority.MinAmount,
		MaxAmount:    s.conf.Priority.MaxAmount,
		TermMonths:   s.conf.Priority.TermMonths,
	}, nil
}

func (s *Service) getCreateConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string) (
	string, error) {
	createConsentResp, err := s.dependencies.RedisAwaiting.GetValue(ctx, consentType, login, bank)
	if err != nil && !errors.Is(err, redisAwaitingProviderModels.ErrRequestIDNotFound) {
		return "", fmt.Errorf("s.dependencies.RedisAwaiting.GetValue: %w", err)
	}

	switch {
	case errors.Is(err, redisAwaitingProviderModels.ErrRequestIDNotFound):
		createConsentResp, err = s.dependencies.RedisApproved.GetValue(ctx, consentType, login, bank)

		switch {
		case errors.Is(err, redisApprovedProviderModels.ErrConsentIDNotFound):
			createConsentResp, err = s.createCreateConsent(ctx, consentType, bank, login)
			if err != nil {
				return "", fmt.Errorf("s.createCreateConsent: %w", err)
			}

			return createConsentResp, nil
		case err != nil:
			return "", fmt.Errorf("s.dependencies.RedisApproved.GetValue: %w", err)
		}

		return createConsentResp, nil
	case err != nil:
		return "", fmt.Errorf("s.dependencies.RedisAwaiting.GetValue: consent_type=%s: %w",
			consentType, err)
	default:
		err = s.dependencies.ConsentService.CheckConsent(ctx, consentServiceModels.CheckConsentRequest{
			ConsentID:   createConsentResp,
			Login:       login,
			Bank:        bank,
			ConsentType: consentType,
		})
		if err != nil {
			return "", fmt.Errorf("s.dependencies.ConsentService.CheckConsent: %w", err)
		}

		err = s.dependencies.RedisAwaiting.Delete(ctx, consentType, login, bank)
		if err != nil {
			return "", fmt.Errorf("s.dependencies.RedisAwaiting.Delete: %w", err)
		}

		resp, err := s.dependencies.OpenApi.GetCreateConsent(ctx, consentType, bank, login, createConsentResp)
		if err != nil {
			return "", fmt.Errorf("s.dependencies.OpenApi.GetConsent: %w", err)
		}

		if resp.ConsentID == "" {
			return "", fmt.Errorf("10 found empty consent id")
		}

		err = s.dependencies.RedisApproved.SetValue(ctx, consentType, login, bank, resp.ConsentID, 0)
		if err != nil {
			return "", fmt.Errorf("s.dependencies.RedisApproved.SetValue: %w", err)
		}

		createConsentResp = resp.ConsentID
	}

	return createConsentResp, nil
}

func (s *Service) getReadConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string) (
	string, error) {
	readConsent, err := s.dependencies.RedisAwaiting.GetValue(ctx, consentType, login, bank)
	if err != nil && !errors.Is(err, redisAwaitingProviderModels.ErrRequestIDNotFound) {
		return "", fmt.Errorf("s.dependencies.RedisAwaiting.GetValue: %w", err)
	}

	switch {
	case errors.Is(err, redisAwaitingProviderModels.ErrRequestIDNotFound):
		readConsent, err = s.dependencies.RedisApproved.GetValue(ctx, consentType, login, bank)

		switch {
		case errors.Is(err, redisApprovedProviderModels.ErrConsentIDNotFound):
			readConsent, err = s.createReadConsent(ctx, consentType, bank, login)
			if err != nil {
				return "", fmt.Errorf("s.createReadConsent: %w", err)
			}

			return readConsent, nil
		case err != nil:
			return "", fmt.Errorf("s.dependencies.RedisApproved.GetValue: %w", err)
		}

		return readConsent, nil
	case err != nil:
		return "", fmt.Errorf("s.dependencies.RedisAwaiting.GetValue: %w", err)
	default:
		err = s.dependencies.ConsentService.CheckConsent(ctx, consentServiceModels.CheckConsentRequest{
			ConsentID:   readConsent,
			Login:       login,
			Bank:        bank,
			ConsentType: consentType,
		})
		if err != nil {
			return "", fmt.Errorf("s.dependencies.ConsentService.CheckConsent: %w", err)
		}

		err = s.dependencies.RedisAwaiting.Delete(ctx, consentType, login, bank)
		if err != nil {
			return "", fmt.Errorf("s.dependencies.RedisAwaiting.Delete: %w", err)
		}

		resp, err := s.dependencies.OpenApi.GetConsent(ctx, consentType, bank, login, readConsent)
		if err != nil {
			return "", fmt.Errorf("s.dependencies.OpenApi.GetConsent: %w", err)
		}

		if resp.Data.ConsentID == "" {
			return "", fmt.Errorf("9 found empty consent id")
		}

		err = s.dependencies.RedisApproved.SetValue(ctx, consentType, login, bank, resp.Data.ConsentID, 0)
		if err != nil {
			return "", fmt.Errorf("s.dependencies.RedisApproved.SetValue: %w", err)
		}

		readConsent = resp.Data.ConsentID

		return readConsent, nil
	}
}

func (s *Service) createCreateConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string,
) (createConsent string, err error) {
	createConsentResp, createConsentErr := s.dependencies.ConsentService.CreateCreateConsent(ctx,
		consentServiceModels.CreateCreateConsentRequest{
			ConsentType: consentType,
			Login:       login,
			BankFrom:    bank,
		})
	if createConsentErr != nil {
		return "", fmt.Errorf("s.dependencies.ConsentService.CreateCreateConsent: %w", createConsentErr)
	}

	if !createConsentResp.AutoApproved {
		if createConsentResp.RequestID == "" {
			return "", fmt.Errorf("6 found empty request id")
		}

		err = s.dependencies.RedisAwaiting.SetValue(ctx, consentType, login, bank, createConsentResp.RequestID,
			0)
		if err != nil {
			return "", fmt.Errorf("s.dependencies.RedisAwaiting.SetValue: %w", err)
		}

		return "", fmt.Errorf("s.dependencies.ConsentService.CreateCreateConsent: %w",
			models.ErrConsentIsNotApproved)
	}

	if createConsentResp.ConsentID == "" {
		return "", fmt.Errorf("5 found empty consent id")
	}

	err = s.dependencies.RedisApproved.SetValue(ctx, consentType, login, bank, createConsentResp.ConsentID,
		consentType.ValidDuration())
	if err != nil {
		return "", fmt.Errorf("s.dependencies.RedisApproved.SetValue: %w", err)
	}

	createConsent = createConsentResp.ConsentID

	return
}

func (s *Service) createReadConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string,
) (readConsent string, err error) {
	createConsentResp, createConsentErr := s.dependencies.OpenApi.CreateReadConsent(ctx, consentType, bank, login)
	if createConsentErr != nil {
		return "", fmt.Errorf("s.dependencies.OpenApi.OpenReadConsent: %w", createConsentErr)
	}

	if !createConsentResp.AutoApproved {
		if createConsentResp.RequestID == "" {
			return "", fmt.Errorf("7 found empty request id")
		}

		err = s.dependencies.RedisAwaiting.SetValue(ctx, consentType, login, bank, createConsentResp.RequestID,
			0)
		if err != nil {
			return "", fmt.Errorf("s.dependencies.RedisAwaiting.SetValue: %w", err)
		}

		return "", fmt.Errorf("s.dependencies.ConsentService.CreateCreateConsent: %w",
			models.ErrConsentIsNotApproved)
	}

	if createConsentResp.ConsentID == "" {
		return "", fmt.Errorf("8 found empty consent id")
	}

	err = s.dependencies.RedisApproved.SetValue(ctx, consentType, login, bank, createConsentResp.ConsentID,
		consentType.ValidDuration())
	if err != nil {
		return "", fmt.Errorf("s.dependencies.RedisApproved.SetValue: %w", err)
	}

	readConsent = createConsentResp.ConsentID

	return
}
