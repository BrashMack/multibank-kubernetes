package approved

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/masharpik/vtb-hack-back/internal/provider/redis/approved/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

type provider struct {
	client            *goredis.Client
	defaultExpiration time.Duration
}

type Provider interface {
	SetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank, consentID string,
		ttl time.Duration) error
	GetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) (string, error)
	Delete(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) error
	Exists(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) (bool, error)
	ExistsByLogin(ctx context.Context, typeConsent enums.ConsentType, login string) (bool, error)
	Close() error
	GetBanksByLogin(ctx context.Context, typeConsent enums.ConsentType, login string) ([]models.ConsentData, error)
}

func Init(ctx context.Context, cfg Config) (Provider, error) {
	opts := &goredis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Client.Host, cfg.Client.Port),
		Password:        cfg.Client.Password,
		DB:              cfg.Client.DB,
		DialTimeout:     cfg.Client.DialTimeout,
		ReadTimeout:     cfg.Client.ReadTimeout,
		WriteTimeout:    cfg.Client.WriteTimeout,
		PoolSize:        cfg.Client.PoolSize,
		MinIdleConns:    cfg.Client.MinIdleConns,
		ConnMaxIdleTime: cfg.Client.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.Client.ConnMaxLifetime,
	}

	if cfg.Client.UseTLS {
		opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := goredis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("client.Ping: %w", err)
	}

	return &provider{
		client:            client,
		defaultExpiration: cfg.DefaultExpiration,
	}, nil
}
