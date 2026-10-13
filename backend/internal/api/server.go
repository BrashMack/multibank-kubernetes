package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/masharpik/vtb-hack-back/docs"
	"github.com/masharpik/vtb-hack-back/internal/api/account"
	"github.com/masharpik/vtb-hack-back/internal/api/card"
	"github.com/masharpik/vtb-hack-back/internal/api/consent"
	"github.com/masharpik/vtb-hack-back/internal/api/deposit"
	"github.com/masharpik/vtb-hack-back/internal/api/user"
	"github.com/masharpik/vtb-hack-back/internal/provider"
	"github.com/masharpik/vtb-hack-back/internal/service"
	"github.com/masharpik/vtb-hack-back/pkg/reflector"
)

type CorsPolitics struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

type Config struct {
	Card               card.Config
	Consent            consent.Config
	User               user.Config
	Account            account.Config
	Deposit            deposit.Config
	CorsPolitics       CorsPolitics
	HiddenHeadersInLog []string
}

type Servers struct {
	Card    *card.Server
	Consent *consent.Server
	User    *user.Server
	Account *account.Server
	Deposit *deposit.Server
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Init(ctx context.Context, cfg Config, services *service.Services, providers *provider.Providers) error {
	var (
		servers    = reflector.DepthNew[Servers]()
		echoServer = echo.New()
		apiV1Group = echoServer.Group("/api/v1")
	)

	echoServer.Validator = &CustomValidator{validator: validator.New()}

	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.CorsPolitics.AllowOrigins,
		AllowMethods:     cfg.CorsPolitics.AllowMethods,
		AllowHeaders:     cfg.CorsPolitics.AllowHeaders,
		AllowCredentials: cfg.CorsPolitics.AllowCredentials,
	}))

	echoServer.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	echoServer.GET("/swagger/*", echoSwagger.WrapHandler)

	if err := servers.Card.Init(ctx, cfg.Card, apiV1Group, card.Dependencies{
		Service:        services.Card,
		JWTProvider:    providers.JWT,
		PaymentService: services.Payment,
	}); err != nil {
		return fmt.Errorf("servers.Card.Init: %w", err)
	}

	if err := servers.Consent.Init(ctx, cfg.Consent, apiV1Group, consent.Dependencies{
		Service:     services.Consent,
		JWTProvider: providers.JWT,
	}); err != nil {
		return fmt.Errorf("servers.Consents.Init: %w", err)
	}

	if err := servers.User.Init(ctx, cfg.User, apiV1Group, user.Dependencies{
		Service: services.User,
	}); err != nil {
		return fmt.Errorf("servers.User.Init: %w", err)
	}

	if err := servers.Account.Init(ctx, cfg.Account, apiV1Group, account.Dependencies{
		Service:     services.Card,
		JWTProvider: providers.JWT,
	}); err != nil {
		return fmt.Errorf("servers.Account.Init: %w", err)
	}

	if err := servers.Deposit.Init(ctx, cfg.Deposit, apiV1Group, deposit.Dependencies{
		JWTProvider: providers.JWT,
		Service:     services.Deposit,
	}); err != nil {
		return fmt.Errorf("servers.Deposit.Init: %w", err)
	}

	echoServer.Use(LogHTTPMiddleware(cfg.HiddenHeadersInLog))

	if err := echoServer.Start(":8080"); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("echoServer.Start: %w", err)
	}

	return nil
}
