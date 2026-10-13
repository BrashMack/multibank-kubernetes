package redis

import (
	"context"
	"fmt"

	"github.com/masharpik/vtb-hack-back/internal/provider/redis/approved"
	"github.com/masharpik/vtb-hack-back/internal/provider/redis/awaiting"
)

type Providers struct {
	Approved approved.Provider
	Awaiting awaiting.Provider
}

type Config struct {
	Approved approved.Config
	Awaiting awaiting.Config
}

func Init(ctx context.Context, cfg Config) (*Providers, error) {
	approvedProvider, err := approved.Init(ctx, cfg.Approved)
	if err != nil {
		return nil, fmt.Errorf("approved.Init: %w", err)
	}

	awaitingProvider, err := awaiting.Init(ctx, cfg.Awaiting)
	if err != nil {
		return nil, fmt.Errorf("awaiting.Init: %w", err)
	}

	return &Providers{
		Approved: approvedProvider,
		Awaiting: awaitingProvider,
	}, nil
}
