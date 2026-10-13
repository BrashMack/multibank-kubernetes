package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type provider struct {
	pool *pgxpool.Pool
}

// Provider - структура методов, которые реализует provider к БД по картам.
type Provider interface {
	ValidateCredentials(ctx context.Context, login, hashedPassword string) error
}

func Init(pool *pgxpool.Pool) Provider {
	return &provider{pool: pool}
}
