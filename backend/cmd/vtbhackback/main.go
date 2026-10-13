// @title       Сервер API мультибанка
// @description Сервер предоставляет API мультибанка, интегрирующего несколько Open Banking API.
// @version     1.0
// @description API сервиса
// @host        localhost:18080
// @BasePath    /
package main

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"

	_ "github.com/masharpik/vtb-hack-back/docs"
	"github.com/masharpik/vtb-hack-back/internal/api"
	"github.com/masharpik/vtb-hack-back/internal/config"
	"github.com/masharpik/vtb-hack-back/internal/provider"
	"github.com/masharpik/vtb-hack-back/internal/service"
	"github.com/masharpik/vtb-hack-back/pkg/slog"
)

func main() {
	ctx := context.Background()

	cfg, loadConfigErr := config.Load()
	if loadConfigErr != nil {
		log.Fatal(loadConfigErr)
	}

	slog.InitLogger(cfg.Logger)

	providers, provInitErr := provider.Init(ctx, cfg.Providers)
	if provInitErr != nil {
		logrus.Fatalf("provider.Init: %s", provInitErr.Error())
	}

	services, serviceInitErr := service.Init(ctx, cfg.Service, providers)
	if serviceInitErr != nil {
		logrus.Fatalf("service.Init: %s", serviceInitErr.Error())
	}

	if apiInitErr := api.Init(ctx, cfg.API, services, providers); apiInitErr != nil {
		logrus.Fatalf("api.Init: %s", apiInitErr.Error())
	}

	return // What a fuck?
}
