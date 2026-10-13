package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/masharpik/vtb-hack-back/internal/provider/storage/card"
	"github.com/masharpik/vtb-hack-back/internal/provider/storage/consent"
	"github.com/masharpik/vtb-hack-back/internal/provider/storage/user"
)

type Entities struct {
	Card    card.Provider
	Consent consent.Provider
	User    user.Provider
}

func Init(ctx context.Context, cfg Config) (entities *Entities, err error) {
	var pool *pgxpool.Pool

	if pool, err = InitPool(ctx, cfg); err != nil {
		return nil, fmt.Errorf("InitPool: %w", err)
	}

	entities, err = makeEntities(ctx, pool, cfg)
	if err != nil {
		return nil, fmt.Errorf("makeEntities: %w", err)
	}

	return entities, nil
}

func InitPool(ctx context.Context, cfg Config) (pool *pgxpool.Pool, err error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pgxCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	pgxCfg.MaxConnLifetime = cfg.MaxConnLifetime
	pgxCfg.MaxConnLifetimeJitter = cfg.MaxConnLifetimeJitter
	pgxCfg.MaxConns = cfg.MaxConns
	pgxCfg.MinConns = cfg.MinConns
	pgxCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	pool, err = pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	return pool, nil
}

func makeEntities(ctx context.Context, pl *pgxpool.Pool, cfg Config) (entities *Entities, err error) {
	entities = &Entities{
		Card:    card.Init(pl),
		Consent: consent.Init(pl),
		User:    user.Init(pl),
	}

	return entities, nil
}
