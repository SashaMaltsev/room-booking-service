package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SashaMaltsev/room-booking-service/internal/app"
	"github.com/SashaMaltsev/room-booking-service/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	api, err := app.NewAPI(ctx, cfg)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}
	defer func() {
		if closeErr := api.Close(); closeErr != nil {
			log.Printf("close app: %v", closeErr)
		}
	}()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("http server listening on %s", api.Server.Addr)
		serverErr <- api.Server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
		defer cancel()

		if err := api.Server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("shutdown server: %v", err)
		}
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("serve http: %v", err)
		}
	}
}
