package consent

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"slices"

	"github.com/masharpik/vtb-hack-back/internal/service/consent/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

func (s *Service) CreateReadConsent(ctx context.Context, req models.CreateReadConsentRequest) (
	*models.CreateReadConsentResponse, error) {
	exists, err := s.dependencies.RedisApproved.Exists(ctx, req.ConsentType, req.Login, req.Bank)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisApproved.Exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("s.dependencies.RedisApproved.Exists: %w", models.ErrConsentAlreadyExists)
	}

	exists, err = s.dependencies.RedisAwaiting.Exists(ctx, req.ConsentType, req.Login, req.Bank)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisAwaiting.Exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("s.dependencies.RedisAwaiting.Exists: %w", models.ErrConsentAlreadyAwaiting)
	}

	resp, err := s.dependencies.OpenApi.CreateReadConsent(ctx, req.ConsentType, req.Bank, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.OpenApi.CreateConsent: %w", err)
	}

	switch {
	case resp.AutoApproved:
		if resp.ConsentID == "" {
			return nil, fmt.Errorf("1 found empty consent id")
		}

		if err = s.dependencies.RedisApproved.SetValue(ctx,
			req.ConsentType, req.Login, req.Bank, resp.ConsentID, 0); err != nil {
			return nil, fmt.Errorf("s.dependencies.RedisApproved.SetValue: %w", err)
		}
	default:
		if resp.RequestID == "" {
			return nil, fmt.Errorf("2 found empty request id")
		}

		if err = s.dependencies.RedisAwaiting.SetValue(ctx,
			req.ConsentType, req.Login, req.Bank, resp.RequestID, 0); err != nil {
			return nil, fmt.Errorf("s.dependencies.RedisAwaiting.SetValue: %w", err)
		}
	}

	return &models.CreateReadConsentResponse{
		AutoApproved: resp.AutoApproved,
	}, nil
}

func (s *Service) CreateWriteConsent(ctx context.Context, req models.CreateWriteConsentRequest) (
	*models.CreateWriteConsentResponse, error) {
	resp, err := s.dependencies.OpenApi.CreateWriteConsent(ctx, req.ConsentType, req.BankFrom, req.Login,
		req.Amount, req.DebtorAccount)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.OpenApi.CreateConsent: %w", err)
	}

	if !resp.AutoApproved {
		return nil, fmt.Errorf("s.dependencies.OpenApi.CreateConsent: %w", models.ErrUnsuppotedNoAutoApproved)
	}

	return &models.CreateWriteConsentResponse{
		RequestID:    resp.RequestID,
		ConsentID:    resp.ConsentID,
		AutoApproved: resp.AutoApproved,
	}, nil
}

func (s *Service) CreateCreateConsent(ctx context.Context, req models.CreateCreateConsentRequest) (
	*models.CreateCreateConsentResponse, error) {
	resp, err := s.dependencies.OpenApi.CreateCreateConsent(ctx, req.ConsentType, req.BankFrom, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.OpenApi.CreateCreateConsent: %w", err)
	}

	return &models.CreateCreateConsentResponse{
		RequestID:    resp.RequestID,
		ConsentID:    resp.ConsentID,
		AutoApproved: resp.AutoApproved,
	}, nil
}

func (s *Service) CheckConsent(ctx context.Context, req models.CheckConsentRequest) error {
	switch req.ConsentType {
	case enums.Create:
		resp, err := s.dependencies.OpenApi.GetCreateConsent(ctx, req.ConsentType, req.Bank, req.Login, req.ConsentID)
		if err != nil {
			return fmt.Errorf("s.dependencies.OpenApi.GetCreateConsent: %w", err)
		}

		if resp.Status != enums.Active {
			return models.ErrConsentNotActive
		}
	case enums.Read:
		resp, err := s.dependencies.OpenApi.GetConsent(ctx, req.ConsentType, req.Bank, req.Login,
			req.ConsentID)
		if err != nil {
			return fmt.Errorf("s.dependencies.OpenApi.GetConsent: consent_type=%s, consent_id=%s, bank=%s: %w",
				req.ConsentType, req.ConsentID, req.Bank, err)
		}

		if !slices.Contains([]enums.ConsentStatus{enums.Approved, enums.Active, enums.Authorized}, resp.Data.Status) {
			return models.ErrConsentNotActive
		}

		return nil
	default:
		return models.ErrUnknownConsentType
	}

	return nil
}

func (s *Service) ListAwaiting(ctx context.Context, req models.ListAwaitingRequest) (
	*models.ListAwaitingResponse, error) {
	banks := make([]models.BankRequest, 0)

	readRequests, err := s.dependencies.RedisAwaiting.GetBanksByLogin(ctx, enums.Read, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisAwaiting.GetBanksByLogin: consent_type=%s: %w",
			enums.Read, err)
	}

	for _, requestData := range readRequests {
		readConsentResp, readConsentErr := s.dependencies.OpenApi.GetConsent(ctx, enums.Read, requestData.Bank,
			req.Login, requestData.RequestID)
		if readConsentErr != nil {
			return nil, fmt.Errorf("s.dependencies.OpenApi.GetConsent: "+
				"consent_type=%s, request_id=%s, bank=%s: %w",
				enums.Read, requestData.RequestID, requestData.Bank, readConsentErr)
		}

		switch readConsentResp.Data.Status {
		case enums.AwaitingAuthorization:
			renamedBankName, ok := s.conf.BankToName[requestData.Bank]
			if !ok {
				renamedBankName = requestData.Bank
			}

			banks = append(banks, models.BankRequest{
				Bank:        renamedBankName,
				RequestType: enums.Read,
				RequestID:   requestData.RequestID,
			})
		case enums.Rejected, enums.Revoked:
			if err = s.dependencies.RedisAwaiting.Delete(ctx, enums.Read, req.Login, requestData.Bank); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisAwaiting.Delete: %w", err)
			}
		case enums.Authorized:
			if err = s.dependencies.RedisAwaiting.Delete(ctx, enums.Read, req.Login, requestData.Bank); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisAwaiting.Delete: %w", err)
			}

			if readConsentResp.Data.ConsentID == "" {
				logrus.Debug("3 found empty consent id")

				continue
			}

			if err = s.dependencies.RedisApproved.SetValue(
				ctx, enums.Read, req.Login, requestData.Bank, readConsentResp.Data.ConsentID, 0); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisApproved.SetValue: %w", err)
			}
		}
	}

	writeRequests, err := s.dependencies.RedisAwaiting.GetBanksByLogin(ctx, enums.Write, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisAwaiting.GetBanksByLogin: consent_type=%s: %w",
			enums.Write, err)
	}

	for _, requestData := range writeRequests {
		renamedBankName, ok := s.conf.BankToName[requestData.Bank]
		if !ok {
			renamedBankName = requestData.Bank
		}

		banks = append(banks, models.BankRequest{
			Bank:        renamedBankName,
			RequestType: enums.Write,
			RequestID:   requestData.RequestID,
		})
	}

	var createBanks = make([]models.BankRequest, 0)

	createRequests, err := s.dependencies.RedisAwaiting.GetBanksByLogin(ctx, enums.Create, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisAwaiting.GetBanksByLogin: consent_type=%s: %w",
			enums.Create, err)
	}

	for _, requestData := range createRequests {
		readConsentResp, readConsentErr := s.dependencies.OpenApi.GetCreateConsent(ctx, enums.Create, requestData.Bank,
			req.Login, requestData.RequestID)
		if readConsentErr != nil {
			return nil, fmt.Errorf("s.dependencies.OpenApi.GetCreateConsent: "+
				"consent_type=%s, request_id=%s, bank=%s: %w",
				enums.Create, requestData.RequestID, requestData.Bank, readConsentErr)
		}

		switch readConsentResp.Status {
		case enums.AwaitingAuthorization:
			renamedBankName, ok := s.conf.BankToName[requestData.Bank]
			if !ok {
				renamedBankName = requestData.Bank
			}

			createBanks = append(createBanks, models.BankRequest{
				Bank:        renamedBankName,
				RequestType: enums.Create,
				RequestID:   requestData.RequestID,
			})
		case enums.Rejected, enums.Revoked:
			if err = s.dependencies.RedisAwaiting.Delete(ctx, enums.Create, req.Login, requestData.Bank); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisAwaiting.Delete: %w", err)
			}
		case enums.Authorized, enums.Approved:
			if err = s.dependencies.RedisAwaiting.Delete(ctx, enums.Create, req.Login, requestData.Bank); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisAwaiting.Delete: %w", err)
			}

			if readConsentResp.ConsentID == "" {
				logrus.Debug("11 found empty consent id")

				continue
			}

			if err = s.dependencies.RedisApproved.SetValue(
				ctx, enums.Create, req.Login, requestData.Bank, readConsentResp.ConsentID, 0); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisApproved.SetValue: %w", err)
			}
		}
	}

	return &models.ListAwaitingResponse{
		BankRequests:       banks,
		BankCreateRequests: createBanks,
	}, nil
}

func (s *Service) ListConfirmed(ctx context.Context, req models.ListConfirmedRequest) (
	*models.ListConfirmedResponse, error) {
	banks := make([]models.BankConsent, 0)

	readConsents, err := s.dependencies.RedisApproved.GetBanksByLogin(ctx, enums.Read, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisApproved.GetBanksByLogin: consent_type=%s: %w",
			enums.Read, err)
	}

	for _, consentData := range readConsents {
		readConsentResp, readConsentErr := s.dependencies.OpenApi.GetConsent(ctx, enums.Read, consentData.Bank,
			req.Login, consentData.ConsentID)
		if readConsentErr != nil {
			return nil, fmt.Errorf("s.dependencies.OpenApi.GetConsent: "+
				"consent_type=%s, consent_id=%s, bank=%s: %w",
				enums.Read, consentData.ConsentID, consentData.Bank, readConsentErr)
		}

		switch readConsentResp.Data.Status {
		case enums.AwaitingAuthorization:
			if err = s.dependencies.RedisApproved.Delete(ctx, enums.Read, req.Login, consentData.Bank); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisApproved.Delete: %w", err)
			}

			if readConsentResp.Data.ConsentID == "" {
				logrus.Debug("4 found empty request id")

				continue
			}

			if err = s.dependencies.RedisAwaiting.SetValue(
				ctx, enums.Read, req.Login, consentData.Bank, readConsentResp.Data.ConsentID, 0); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisAwaiting.SetValue: %w", err)
			}
		case enums.Rejected, enums.Revoked:
			if err = s.dependencies.RedisApproved.Delete(ctx, enums.Read, req.Login, consentData.Bank); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisApproved.Delete: %w", err)
			}
		case enums.Authorized:
			renamedBankName, ok := s.conf.BankToName[consentData.Bank]
			if !ok {
				renamedBankName = consentData.Bank
			}

			banks = append(banks, models.BankConsent{
				Bank:        renamedBankName,
				ConsentType: enums.Read,
				ConsentID:   consentData.ConsentID,
			})
		}
	}

	writeConsents, err := s.dependencies.RedisApproved.GetBanksByLogin(ctx, enums.Write, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisApproved.GetBanksByLogin: consent_type=%s: %w",
			enums.Write, err)
	}

	for _, consentData := range writeConsents {
		renamedBankName, ok := s.conf.BankToName[consentData.Bank]
		if !ok {
			renamedBankName = consentData.Bank
		}

		banks = append(banks, models.BankConsent{
			Bank:        renamedBankName,
			ConsentType: enums.Write,
			ConsentID:   consentData.ConsentID,
		})
	}

	var createBanks = make([]models.BankConsent, 0)

	createConsents, err := s.dependencies.RedisApproved.GetBanksByLogin(ctx, enums.Create, req.Login)
	if err != nil {
		return nil, fmt.Errorf("s.dependencies.RedisApproved.GetBanksByLogin: consent_type=%s: %w",
			enums.Create, err)
	}

	for _, consentData := range createConsents {
		readConsentResp, readConsentErr := s.dependencies.OpenApi.GetCreateConsent(ctx, enums.Create, consentData.Bank,
			req.Login, consentData.ConsentID)
		if readConsentErr != nil {
			return nil, fmt.Errorf("s.dependencies.OpenApi.GetCreateConsent: "+
				"consent_type=%s, consent_id=%s, bank=%s: %w",
				enums.Create, consentData.ConsentID, consentData.Bank, readConsentErr)
		}

		switch readConsentResp.Status {
		case enums.AwaitingAuthorization:
			if err = s.dependencies.RedisApproved.Delete(ctx, enums.Create, req.Login, consentData.Bank); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisApproved.Delete: %w", err)
			}

			if readConsentResp.ConsentID == "" {
				logrus.Debug("12 found empty consent id")

				continue
			}

			if err = s.dependencies.RedisAwaiting.SetValue(
				ctx, enums.Create, req.Login, consentData.Bank, readConsentResp.ConsentID, 0); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisAwaiting.SetValue: %w", err)
			}
		case enums.Rejected, enums.Revoked:
			if err = s.dependencies.RedisApproved.Delete(ctx, enums.Create, req.Login, consentData.Bank); err != nil {
				return nil, fmt.Errorf("s.dependencies.RedisApproved.Delete: %w", err)
			}
		case enums.Authorized, enums.Approved:
			renamedBankName, ok := s.conf.BankToName[consentData.Bank]
			if !ok {
				renamedBankName = consentData.Bank
			}

			createBanks = append(createBanks, models.BankConsent{
				Bank:        renamedBankName,
				ConsentType: enums.Write,
				ConsentID:   consentData.ConsentID,
			})
		}
	}

	return &models.ListConfirmedResponse{
		BankConsents:       banks,
		BankCreateConsents: createBanks,
	}, nil
}

func (s *Service) ListActive(ctx context.Context, req models.ListActiveRequest) (
	*models.ListActiveResponse, error) {
	awaitingBanks, err := s.ListAwaiting(ctx, models.ListAwaitingRequest{
		Login: req.Login,
	})
	if err != nil {
		return nil, fmt.Errorf("s.ListAwaiting: %w", err)
	}

	confirmedBanks, err := s.ListConfirmed(ctx, models.ListConfirmedRequest{
		Login: req.Login,
	})
	if err != nil {
		return nil, fmt.Errorf("s.ListConfirmed: %w", err)
	}

	return &models.ListActiveResponse{
		ConfirmedBanks:       confirmedBanks.BankConsents,
		ConfirmedCreateBanks: confirmedBanks.BankCreateConsents,
		AwaitingBanks:        awaitingBanks.BankRequests,
		AwaitingCreateBanks:  awaitingBanks.BankCreateRequests,
	}, nil
}
