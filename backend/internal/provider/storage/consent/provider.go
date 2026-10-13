package consent

import "github.com/jackc/pgx/v5/pgxpool"

type provider struct {
	pool *pgxpool.Pool
}

// Provider - структура методов, которые реализует provider к БД по картам.
type Provider interface {
}

func Init(pool *pgxpool.Pool) Provider {
	return &provider{pool: pool}
}
