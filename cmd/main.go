package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sso-bff/internal/config"
	"sso-bff/internal/db"
	"sso-bff/internal/lib"
	"sso-bff/modules"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env only for local dev
	if os.Getenv("ENV") != "prod" {
		_ = godotenv.Load()
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	resources, err := db.NewResources(ctx, cfg.DB.DSN, cfg.Store.StoreAddress, cfg.Store.StorePassword, cfg.Store.StoreDB)
	if err != nil {
		log.Fatalf("failed to initialize resources: %v", err)
	}
	defer resources.Close()

	sdks := modules.NewSDKs(cfg)

	app := &lib.App{
		Config:    cfg,
		Resources: resources,
		SDK:       sdks,
	}

	_ = app
}
