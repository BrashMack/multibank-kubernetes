package config

import (
	"fmt"
	"net"
	"net/url"

	"github.com/masharpik/vtb-hack-back/internal/api"
	"github.com/masharpik/vtb-hack-back/internal/provider"
	"github.com/masharpik/vtb-hack-back/internal/service"
	"github.com/masharpik/vtb-hack-back/pkg/loadconfig"
	"github.com/masharpik/vtb-hack-back/pkg/slog"
)

type AppConfig struct {
	API       api.Config
	Service   service.Config
	Providers provider.Config
	Logger    slog.Config
}

// Load - метод парсинга файла кофигурации проекта в структуру конфига.
func Load() (*AppConfig, error) {
	var cfg = new(AppConfig)

	if err := loadconfig.To(cfg); err != nil {
		return nil, fmt.Errorf("loadconfig.To: %w", err)
	}

	if len(cfg.Providers.Storage.Host) != 0 {
		cfg.Providers.Storage.DSN = makeDSN(
			"postgres",
			cfg.Providers.Storage.User, cfg.Providers.Storage.Password,
			cfg.Providers.Storage.Host, cfg.Providers.Storage.Port,
			cfg.Providers.Storage.Database,
		)
	}

	return cfg, nil
}

func makeDSN(schema, user, password, host, port, database string) string {
	dsn := new(url.URL)
	dsn.Scheme = schema
	dsn.User = url.UserPassword(user, password)
	dsn.Host = net.JoinHostPort(host, port)
	dsn = dsn.JoinPath(database)
	values := dsn.Query()
	dsn.RawQuery = values.Encode()

	return dsn.String()
}
