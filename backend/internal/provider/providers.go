package provider

import (
	"context"
	"fmt"

	"github.com/masharpik/vtb-hack-back/internal/provider/jwt"
	"github.com/masharpik/vtb-hack-back/internal/provider/openapi"
	"github.com/masharpik/vtb-hack-back/internal/provider/redis"
	"github.com/masharpik/vtb-hack-back/internal/provider/storage"
)

type Providers struct {
	Entities *storage.Entities
	OpenAPI  *openapi.Client
	JWT      *jwt.Provider
	Redis    *redis.Providers
}

type Config struct {
	Storage storage.Config
	OpenAPI openapi.Config
	JWT     jwt.Config
	Redis   redis.Config
}

func Init(ctx context.Context, cfg Config) (providers *Providers, err error) {
	providers = new(Providers)

	providers.Entities, err = storage.Init(ctx, cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("storage.Init: %w", err)
	}

	providers.OpenAPI, err = openapi.Init(ctx, cfg.OpenAPI)
	if err != nil {
		return nil, fmt.Errorf("openapi.Init: %w", err)
	}

	providers.JWT, err = jwt.Init(ctx, cfg.JWT)
	if err != nil {
		return nil, fmt.Errorf("jwt.Init: %w", err)
	}

	providers.Redis, err = redis.Init(ctx, cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("redis.Init: %w", err)
	}

	return providers, nil
}
