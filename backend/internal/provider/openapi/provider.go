package openapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

type Client struct {
	client *http.Client
	// tokenForBank - мапа для получения данных токена для банка.
	tokenForBank map[enums.Bank]models.TokenData
	// baseUrlForBank - мапа для получения url'а для банка.
	baseUrlForBank map[enums.Bank]string
	// bankForHostname - мапа для получения названия банка для домена.
	bankForHostname map[string]enums.Bank
	bankToCode      map[enums.Bank]string
	teamID          string
	teamName        string
	teamSecretID    string
	mu              *sync.Mutex
}

type Provider interface {
	CreateReadConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string) (
		*models.CreateReadConsentResponse, error)
	CreateWriteConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login, amount,
		debtorAccount string) (*models.CreateWriteConsentResponse, error)
	CreateCreateConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string) (
		*models.CreateCreateConsentResponse, error)
	GetConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank, login string,
		consentID string) (*models.GetConsentResponse, error)
	GetAccounts(ctx context.Context, bank enums.Bank, login, consentID string) (*models.GetAccountsResponse, error)
	GetAccountBalance(ctx context.Context, bank enums.Bank, login, consentID string) (
		*models.GetAccountBalanceResponse, error)
	GetAccountTransactions(ctx context.Context, bank enums.Bank, consentID, accountID, login string, page, limit int) (
		*models.GetAccountTransactionsResponse, error)
	CreateAgreement(ctx context.Context, bank enums.Bank, login, consentID, productID,
		srcAccountNumber string, amount float64, termMonths int) (*models.CreateAgreementResponse, error)
	GetAgreement(ctx context.Context, bank enums.Bank, login, consentID, agreementID string) (
		*models.GetAgreementResponse, error)
	GetCreateConsent(ctx context.Context, consentType enums.ConsentType, bank enums.Bank,
		login, consentID string) (*models.GetCreateConsentResponse, error)
	GetAgreements(ctx context.Context, bank enums.Bank, login, consentID string) (
		*models.GetAgreementsResponse, error)
	GetProducts(ctx context.Context, bank enums.Bank, productType enums.AgreementType, consentID string) (
		*models.GetProductsResponse, error)
}

func Init(ctx context.Context, cfg Config) (*Client, error) {
	client := &Client{
		tokenForBank:    make(map[enums.Bank]models.TokenData),
		bankForHostname: make(map[string]enums.Bank),
		baseUrlForBank:  cfg.BaseUrls,
		teamID:          cfg.TeamID,
		teamSecretID:    cfg.TeamSecretID,
		bankToCode:      cfg.BankToCode,
		teamName:        fmt.Sprintf("Multibank_%s", cfg.TeamID),
		mu:              &sync.Mutex{},
	}

	var err error

	if client.client, err = client.connect(ctx, cfg); err != nil {
		return nil, fmt.Errorf("client.connect: %w", err)
	}

	if err = client.initBankForHostnames(); err != nil {
		return nil, fmt.Errorf("client.initBankForHostnames: %w", err)
	}

	if cfg.Disabled {
		return client, nil
	}

	if err = client.initTokenForBanks(); err != nil {
		return nil, fmt.Errorf("client.initTokenForBanks: %w", err)
	}

	return client, nil
}

// Connect - метод создания http-клиента для запросов в открытые банковские API.
func (c *Client) connect(ctx context.Context, cfg Config) (*http.Client, error) {
	var tlsConfig *tls.Config

	switch {
	case cfg.Client.Insecure:
		tlsConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		}
	default:
		cert, err := tls.LoadX509KeyPair(cfg.Client.ClientCert, cfg.Client.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("tls.LoadX509KeyPair: %w", err)
		}

		caCert, err := os.ReadFile(cfg.Client.RootCACert)
		if err != nil {
			return nil, fmt.Errorf("os.ReadFile: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("appendCertsFromPEM")
		}

		tlsConfig = &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
	}

	baseTransport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		MaxIdleConns:        cfg.Client.MaxIdleConns,
		MaxIdleConnsPerHost: cfg.Client.MaxIdleConnsPerHost,
		IdleConnTimeout:     cfg.Client.IdleConnTimeout,
	}

	transport := http.RoundTripper(baseTransport)

	transport = c.authTransport(transport)
	transport = c.logTransport(transport, cfg.Client.HiddenHeadersInLog)
	transport = c.retryTransport(transport, cfg.Client.Retry.Codes, cfg.Client.Retry.Period, cfg.Client.Retry.Deadline)

	return &http.Client{
		Transport: transport,
		Timeout:   cfg.Client.Timeout,
	}, nil
}

// InitTokenForBanks - заполняет токен для каждого банка.
func (c *Client) initTokenForBanks() error {
	for bank, baseUrl := range c.baseUrlForBank {
		accessToken, err := c.getToken(baseUrl)
		if err != nil {
			return fmt.Errorf("c.getToken: %w", err)
		}

		c.tokenForBank[bank] = accessToken
	}

	return nil
}

// initBankForHostnames - заполняет мапу банка для его хоста.
func (c *Client) initBankForHostnames() error {
	for bank, baseUrl := range c.baseUrlForBank {
		u, err := url.Parse(baseUrl)
		if err != nil {
			return fmt.Errorf("url.Parse for bank %s: %w", bank, err)
		}

		c.bankForHostname[u.Hostname()] = bank
	}

	return nil
}
