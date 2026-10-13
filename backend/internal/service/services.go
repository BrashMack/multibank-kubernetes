package service

import (
	"context"
	"fmt"

	"github.com/masharpik/vtb-hack-back/internal/provider"
	"github.com/masharpik/vtb-hack-back/internal/service/card"
	"github.com/masharpik/vtb-hack-back/internal/service/consent"
	"github.com/masharpik/vtb-hack-back/internal/service/deposit"
	"github.com/masharpik/vtb-hack-back/internal/service/payment"
	"github.com/masharpik/vtb-hack-back/internal/service/user"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
	"github.com/masharpik/vtb-hack-back/pkg/reflector"
)

type Config struct {
	BankToName map[enums.Bank]enums.Bank
	Card       card.Config
	Consent    consent.Config
	User       user.Config
	Payment    payment.Config
	Deposit    deposit.Config
}

type Services struct {
	Card    *card.Service
	Payment *payment.Service
	Consent *consent.Service
	User    *user.Service
	Deposit *deposit.Service
}

func Init(ctx context.Context, cfg Config, providers *provider.Providers) (*Services, error) {
	var services = reflector.DepthNew[Services]()

	if err := services.init(ctx, cfg, providers); err != nil {
		return nil, fmt.Errorf("services.init: %w", err)
	}

	if err := services.start(ctx); err != nil {
		return nil, fmt.Errorf("services.start: %w", err)
	}

	return services, nil
}

func (s *Services) init(ctx context.Context, cfg Config, providers *provider.Providers) (err error) {
	cfg.Card.BankToName = cfg.BankToName
	cfg.Consent.BankToName = cfg.BankToName
	cfg.Deposit.BankToName = cfg.BankToName

	banks := make([]enums.Bank, 0, len(cfg.BankToName))
	for bankName, _ := range cfg.BankToName {
		banks = append(banks, bankName)
	}

	banksWithoutVBank := make([]enums.Bank, 0, len(cfg.BankToName))
	for bankName, _ := range cfg.BankToName {
		if bankName == enums.VirtualBank {
			continue
		}

		banksWithoutVBank = append(banksWithoutVBank, bankName)
	}

	cfg.Deposit.Banks = banksWithoutVBank
	cfg.Payment.Banks = banks

	if err = s.Consent.Init(cfg.Consent, consent.Dependencies{
		OpenApi:       providers.OpenAPI,
		RedisApproved: providers.Redis.Approved,
		RedisAwaiting: providers.Redis.Awaiting,
	}); err != nil {
		return fmt.Errorf("s.Consent.Init: %w", err)
	}

	if err = s.Card.Init(cfg.Card, card.Dependencies{
		OpenApi:        providers.OpenAPI,
		ConsentService: s.Consent,
		RedisApproved:  providers.Redis.Approved,
	}); err != nil {
		return fmt.Errorf("s.Card.Init: %w", err)
	}

	if err = s.Payment.Init(cfg.Payment, payment.Dependencies{
		OpenApi:        providers.OpenAPI,
		ConsentService: s.Consent,
		CardService:    s.Card,
	}); err != nil {
		return fmt.Errorf("s.Payment.Init: %w", err)
	}

	if err = s.User.Init(cfg.User, user.Dependencies{
		Provider:    providers.Entities.User,
		JWTProvider: providers.JWT,
	}); err != nil {
		return fmt.Errorf("s.User.Init: %w", err)
	}

	if err = s.Deposit.Init(cfg.Deposit, deposit.Dependencies{
		OpenApi:        providers.OpenAPI,
		ConsentService: s.Consent,
		RedisApproved:  providers.Redis.Approved,
		RedisAwaiting:  providers.Redis.Awaiting,
	}); err != nil {
		return fmt.Errorf("s.Deposit.Init: %w", err)
	}

	return nil
}

// Start - метод запускает сервисные методы Start для фоновой логики.
func (s *Services) start(ctx context.Context) error {
	return nil
}
