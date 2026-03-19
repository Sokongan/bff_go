package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"sso-bff/internal/config"
	"sso-bff/internal/db"
	"sso-bff/internal/factory"
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

	module, err := factory.NewHandlers(cfg, resources, sdks)
	if err != nil {
		log.Fatalf("failed to build modules: %v", err)
	}

	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: module.Handler,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	log.Printf("starting HTTP server on %s", cfg.ServerAddr)

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	case <-ctx.Done():
		if shutdownErr := server.Shutdown(context.Background()); shutdownErr != nil {
			log.Fatalf("shutdown error: %v", shutdownErr)
		}
		if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}
}
