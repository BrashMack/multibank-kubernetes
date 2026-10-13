package jwt

import (
	"context"
)

type Provider struct {
	conf Config
}

func Init(ctx context.Context, cfg Config) (*Provider, error) {
	return &Provider{
		conf: cfg,
	}, nil
}
