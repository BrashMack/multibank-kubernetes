package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

func (c *Client) CreateReadConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string) (
	*models.CreateReadConsentResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bank)
	}

	switch consentType {
	case enums.Read:
		url := fmt.Sprintf("%s/account-consents/request", baseUrl)

		reqBody := models.CreateReadConsentRequest{
			ClientID:           login,
			Permissions:        consentType.Permissions(),
			Reason:             consentType.Reason(),
			RequestingBank:     c.teamID,
			RequestingBankName: c.teamName,
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Requesting-Bank", c.teamID)

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("c.client.Do: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBytes, _ := io.ReadAll(resp.Body)

			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
				resp.StatusCode, string(respBytes))
		}

		var result = new(models.CreateReadConsentResponse)
		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return nil, fmt.Errorf("json.Decode: %w", err)
		}

		if result.ConsentID == "" && result.AutoApproved {
			return nil, models.ErrEmptyConsentID
		}

		if result.RequestID == "" {
			return nil, models.ErrEmptyRequestID
		}

		return result, nil
	default:
		return nil, models.ErrUnknownConsentType
	}
}

func (c *Client) CreateWriteConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login, amount,
	debtorAccount string) (
	*models.CreateWriteConsentResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bank)
	}

	switch consentType {
	case enums.Write:
		url := fmt.Sprintf("%s/payment-consents/request", baseUrl)

		floatAmount, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return nil, fmt.Errorf("strconv.ParseFloat: %w: %w", err, models.ErrInvalidAmount)
		}

		reqBody := models.CreateWriteConsentRequest{
			RequestingBank: c.teamID,
			ClientID:       login,
			ConsentType:    "single_use",
			DebtorAccount:  debtorAccount,
			Amount:         floatAmount,
			Reference:      "Оплата услуг",
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Requesting-Bank", c.teamID)

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("c.client.Do: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBytes, _ := io.ReadAll(resp.Body)

			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
				resp.StatusCode, string(respBytes))
		}

		var result = new(models.CreateWriteConsentResponse)
		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return nil, fmt.Errorf("json.Decode: %w", err)
		}

		if result.ConsentID == "" && result.AutoApproved {
			return nil, models.ErrEmptyConsentID
		}

		if result.RequestID == "" {
			return nil, models.ErrEmptyRequestID
		}

		return result, nil
	default:
		return nil, models.ErrUnknownConsentType
	}
}

func (c *Client) CreateCreateConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string,
) (*models.CreateCreateConsentResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bank)
	}

	switch consentType {
	case enums.Create:
		url := fmt.Sprintf("%s/product-agreement-consents/request", baseUrl)

		reqBody := models.CreateCreateConsentRequest{
			RequestingBank:         c.teamID,
			ClientID:               login,
			ReadProductAgreements:  consentType.CanReadProducts(),
			OpenProductAgreements:  consentType.CanOpenProducts(),
			CloseProductAgreements: consentType.CanCloseProducts(),
			MaxAmount:              consentType.MaxAmount(),
			ValidUntil:             consentType.ValidUntil(),
			Reason:                 consentType.Reason(),
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		q := req.URL.Query()
		q.Set("client_id", login)
		req.URL.RawQuery = q.Encode()

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("c.client.Do: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBytes, _ := io.ReadAll(resp.Body)

			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
				resp.StatusCode, string(respBytes))
		}

		var result = new(models.CreateCreateConsentResponse)
		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return nil, fmt.Errorf("json.Decode: %w", err)
		}

		if result.ConsentID == "" && result.AutoApproved {
			return nil, models.ErrEmptyConsentID
		}

		if result.RequestID == "" {
			return nil, models.ErrEmptyRequestID
		}

		return result, nil
	default:
		return nil, models.ErrUnknownConsentType
	}
}

func (c *Client) GetCreateConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank,
	login, consentID string) (*models.GetCreateConsentResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bank)
	}

	switch consentType {
	case enums.Create:
		url := fmt.Sprintf("%s/product-agreement-consents/%s", baseUrl, consentID)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		q := req.URL.Query()
		q.Set("client_id", login)
		req.URL.RawQuery = q.Encode()

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("c.client.Do: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBytes, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusNotFound {
				return nil, fmt.Errorf("c.client.Get: status_code=%d, response_body=%s: %w",
					resp.StatusCode, string(respBytes), models.ErrConsentIDNotFound)
			}

			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
				resp.StatusCode, string(respBytes))
		}

		var result = new(models.GetCreateConsentResponse)
		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return nil, fmt.Errorf("json.Decode: %w", err)
		}

		return result, nil
	default:
		return nil, models.ErrUnknownConsentType
	}
}

func (c *Client) CreateAgreement(ctx context.Context, bank enums.Bank, login, consentID, productID,
	srcAccountNumber string, amount float64, termMonths int) (*models.CreateAgreementResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bank)
	}

	url := fmt.Sprintf("%s/product-agreements", baseUrl)

	reqBody := models.CreateAgreementRequest{
		ProductID:       productID,
		Amount:          amount,
		TermMonths:      termMonths,
		SourceAccountID: srcAccountNumber,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Set("client_id", login)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
			resp.StatusCode, string(respBytes))
	}

	var result = new(models.CreateAgreementResponse)
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	if result.Data.AgreementID == "" {
		return nil, models.ErrEmptyAgreementID
	}

	return result, nil
}

func (c *Client) GetAgreement(ctx context.Context, bank enums.Bank, login, consentID, agreementID string) (
	*models.GetAgreementResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bank)
	}

	url := fmt.Sprintf("%s/product-agreements/%s", baseUrl, agreementID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("X-Requesting-Bank", c.teamID)
	req.Header.Set("X-Product-Agreement-Consent-Id", consentID)

	q := req.URL.Query()
	q.Set("client_id", login)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
			resp.StatusCode, string(respBytes))
	}

	var result = new(models.GetAgreementResponse)
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	if result.Data.AgreementID == "" {
		return nil, models.ErrEmptyAgreementID
	}

	return result, nil
}

func (c *Client) GetAgreements(ctx context.Context, bank enums.Bank, login, consentID string) (
	*models.GetAgreementsResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bank)
	}

	url := fmt.Sprintf("%s/product-agreements", baseUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("X-Requesting-Bank", c.teamID)
	req.Header.Set("X-Product-Agreement-Consent-Id", consentID)

	q := req.URL.Query()
	q.Set("client_id", login)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
			resp.StatusCode, string(respBytes))
	}

	var result = new(models.GetAgreementsResponse)
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	return result, nil
}

func (c *Client) GetProducts(ctx context.Context, bank enums.Bank, productType enums.AgreementType, consentID string) (
	*models.GetProductsResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bank)
	}

	url := fmt.Sprintf("%s/products", baseUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("X-Requesting-Bank", c.teamID)
	req.Header.Set("X--Consent-Id", consentID)

	q := req.URL.Query()
	q.Set("product_type", string(productType))
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
			resp.StatusCode, string(respBytes))
	}

	var result = new(models.GetProductsResponse)
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	return result, nil
}

func (c *Client) CreatePayment(ctx context.Context, consentID string, bankFrom enums.Bank, bankTo enums.Bank, amount,
	currency, login, accountIDFrom, accountIDTo, schemeNameFrom, schemeNameTo string) (
	*models.CreatePaymentResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bankFrom]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bankFrom)
	}

	url := fmt.Sprintf("%s/payments", baseUrl)

	var req *http.Request

	if bankFrom == bankTo {
		reqBody := models.CreatePaymentInRequest{
			Data: models.Data{
				Initiation: models.Initiation{
					InstructedAmount: models.InstructedAmount{
						Amount:   amount,
						Currency: currency,
					},
					DebtorAccount: models.DebtorAccount{
						SchemeName:     schemeNameFrom,
						Identification: accountIDFrom,
					},
					CreditorAccount: models.CreditorAccount{
						SchemeName:     schemeNameTo,
						Identification: accountIDTo,
					},
				},
			},
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %w", err)
		}

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
		}
	} else {
		var bankCode string

		bankCode, ok = c.bankToCode[bankTo]
		if !ok {
			bankCode = string(bankTo)
		}

		reqBody := models.CreatePaymentOutRequest{
			Data: models.DataWithDescBank{
				Initiation: models.InitiationWithDescBank{
					InstructedAmount: models.InstructedAmount{
						Amount:   amount,
						Currency: currency,
					},
					DebtorAccount: models.DebtorAccount{
						SchemeName:     "RU.CBR.PAN",
						Identification: accountIDFrom,
					},
					CreditorAccount: models.CreditorAccountWithBank{
						SchemeName:     "RU.CBR.PAN",
						Identification: accountIDTo,
						BankCode:       bankCode,
					},
				},
			},
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %w", err)
		}

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requesting-Bank", c.teamID)
	req.Header.Set("X-Payment-Consent-Id", consentID)

	q := req.URL.Query()
	q.Set("client_id", login)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBytes, _ := io.ReadAll(resp.Body)

		if strings.Contains(string(respBytes), "Insufficient funds") {
			return nil, models.ErrFailedPaymentStatus
		}

		return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
			resp.StatusCode, string(respBytes))
	}

	var result = new(models.CreatePaymentResponse)
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	result.Data.Status = result.Data.Status.ToLowerTrimSpace()
	if !slices.Contains(enums.SupportedPaymentStatus, result.Data.Status) {
		return nil, models.ErrUnknownPaymentStatus
	}

	if result.Data.Status == enums.FailedPaymentStatus {
		return nil, models.ErrFailedPaymentStatus
	}

	return result, nil
}

func (c *Client) GetConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string,
	consentID string) (*models.GetConsentResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not in c.baseUrls", bank)
	}

	switch consentType {
	case enums.Read:
		url := fmt.Sprintf("%s/account-consents/%s", baseUrl, consentID)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
		}

		req.Header.Set("x-fapi-interaction-id", login)

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("c.client.Do: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBytes, _ := io.ReadAll(resp.Body)

			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s",
				resp.StatusCode, string(respBytes))
		}

		var result = new(models.GetConsentResponse)
		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return nil, fmt.Errorf("json.Decode: %w", err)
		}

		result.Data.Status = result.Data.Status.ToLowerTrimSpace()

		if !slices.Contains(enums.SupportedConsentStatuses, result.Data.Status) {
			return nil, models.ErrUnknownConsentStatus
		}

		return result, nil
	default:
		return nil, models.ErrUnknownConsentType
	}
}

func (c *Client) GetAccounts(ctx context.Context, bank enums.Bank, login, consentID string) (
	*models.GetAccountsResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not found", bank)
	}

	url := fmt.Sprintf("%s/accounts", baseUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("X-Requesting-Bank", c.teamID)
	req.Header.Set("X-Consent-Id", consentID)

	q := req.URL.Query()
	q.Set("client_id", login)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("get accounts: status=%d body=%s", resp.StatusCode, string(respBytes))
	}

	var result = new(models.GetAccountsResponse)

	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	for i, acc := range result.Data.Accounts {
		acc.Status = acc.Status.ToLowerTrimSpace()

		if !slices.Contains(enums.SupportedAccountStatuses, acc.Status) {
			return nil, models.ErrUnknownAccountStatus
		}

		result.Data.Accounts[i].Status = acc.Status

		acc.AccountSubType = acc.AccountSubType.ToLowerTrimSpace()

		if !slices.Contains(enums.SupportedAccountTypes, acc.AccountSubType) {
			return nil, models.ErrUnknownAccountType
		}

		result.Data.Accounts[i].AccountSubType = acc.AccountSubType
	}

	return result, nil
}

func (c *Client) GetAccount(ctx context.Context, bank enums.Bank, login, consentID, accountID string) (
	*models.GetAccountResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not found", bank)
	}

	url := fmt.Sprintf("%s/accounts/%s", baseUrl, accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("X-Requesting-Bank", c.teamID)
	req.Header.Set("X-Consent-Id", consentID)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)

		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrAccountNotFound)
		case http.StatusForbidden:
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("get accounts: status=%d body=%s", resp.StatusCode, string(respBytes))
	}

	var result = new(models.GetAccountResponse)

	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	for i, acc := range result.Data.Accounts {
		acc.Status = acc.Status.ToLowerTrimSpace()

		if !slices.Contains(enums.SupportedAccountStatuses, acc.Status) {
			return nil, models.ErrUnknownAccountStatus
		}

		result.Data.Accounts[i].Status = acc.Status

		acc.AccountSubType = acc.AccountSubType.ToLowerTrimSpace()

		if !slices.Contains(enums.SupportedAccountTypes, acc.AccountSubType) {
			return nil, models.ErrUnknownAccountType
		}

		result.Data.Accounts[i].AccountSubType = acc.AccountSubType
	}

	return result, nil
}

func (c *Client) GetCards(ctx context.Context, bank enums.Bank, login, consentID string) (
	*models.GetCardsResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not found", bank)
	}

	url := fmt.Sprintf("%s/cards", baseUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("X-Requesting-Bank", c.teamID)
	req.Header.Set("X-Consent-Id", consentID)

	q := req.URL.Query()
	q.Set("client_id", login)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("c.client.Do: status=%d body=%s", resp.StatusCode, string(respBytes))
	}

	var result = new(models.GetCardsResponse)

	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	for i, currCard := range result.Data.Cards {
		currCard.Status = currCard.Status.ToLowerTrimSpace()

		if !slices.Contains(enums.SupportedCardStatuses, currCard.Status) {
			return nil, models.ErrUnknownCardStatus
		}

		result.Data.Cards[i].Status = currCard.Status

		currCard.CardType = currCard.CardType.ToLowerTrimSpace()

		if !slices.Contains(enums.SupportedCardTypes, currCard.CardType) {
			return nil, models.ErrUnknownCardType
		}

		result.Data.Cards[i].CardType = currCard.CardType
	}

	return result, nil
}

func (c *Client) GetCard(ctx context.Context, bank enums.Bank, login, consentID, cardID string) (
	*models.GetCardResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not found", bank)
	}

	url := fmt.Sprintf("%s/cards/%s", baseUrl, cardID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("X-Requesting-Bank", c.teamID)
	req.Header.Set("X-Consent-Id", consentID)
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Set("client_id", login)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)

		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrCardNotFound)
		case http.StatusForbidden:
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("c.client.Do: status=%d body=%s", resp.StatusCode, string(respBytes))
	}

	var result = new(models.GetCardResponse)

	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	result.CardData.Status = result.CardData.Status.ToLowerTrimSpace()

	if !slices.Contains(enums.SupportedCardStatuses, result.CardData.Status) {
		return nil, models.ErrUnknownCardStatus
	}

	result.CardData.CardType = result.CardData.CardType.ToLowerTrimSpace()

	if !slices.Contains(enums.SupportedCardTypes, result.CardData.CardType) {
		return nil, models.ErrUnknownCardType
	}

	return result, nil
}

func (c *Client) GetAccountBalance(ctx context.Context, bank enums.Bank, consentID, accountID string) (
	*models.GetAccountBalanceResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not found", bank)
	}

	url := fmt.Sprintf("%s/accounts/%s/balances", baseUrl, accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("X-Consent-Id", consentID)
	req.Header.Set("X-Requesting-Bank", c.teamID)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("c.client.Do: status_code=%d, response_body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("c.client.Do: status_code=%d response_body=%s",
			resp.StatusCode, string(respBytes))
	}

	var result = new(models.GetAccountBalanceResponse)
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	for i, balance := range result.Data.Balance {
		balance.Type = balance.Type.ToLowerTrimSpace()

		if !slices.Contains(enums.SupportedBalanceTypes, balance.Type) {
			return nil, models.ErrUnknownBalanceType
		}

		result.Data.Balance[i].Type = balance.Type.ToLowerTrimSpace()
	}

	return result, nil
}

func (c *Client) GetAccountTransactions(ctx context.Context, bank enums.Bank, consentID, accountID, login string,
	page, limit int) (*models.GetAccountTransactionsResponse, error) {
	baseUrl, ok := c.baseUrlForBank[bank]
	if !ok {
		return nil, fmt.Errorf("bank %s not found", bank)
	}

	url := fmt.Sprintf("%s/accounts/%s/transactions", baseUrl, accountID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("X-Consent-Id", consentID)
	req.Header.Set("X-Requesting-Bank", c.teamID)

	q := req.URL.Query()
	if page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	q.Set("client_id", login)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("c.client.Do: status=%d body=%s: %w",
				resp.StatusCode, string(respBytes), models.ErrConsentForbidden)
		}

		return nil, fmt.Errorf("c.client.Do: status=%d body=%s", resp.StatusCode, string(respBytes))
	}

	var result = new(models.GetAccountTransactionsResponse)
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	return result, nil
}
