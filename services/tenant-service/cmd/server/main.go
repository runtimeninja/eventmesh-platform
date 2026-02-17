package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/runtimeninja/eventmesh-platform/services/tenant-service/internal/config"
	"github.com/runtimeninja/eventmesh-platform/services/tenant-service/internal/db"
	apphttp "github.com/runtimeninja/eventmesh-platform/services/tenant-service/internal/http"
	"github.com/runtimeninja/eventmesh-platform/services/tenant-service/internal/observability"
	"github.com/runtimeninja/eventmesh-platform/services/tenant-service/internal/tenants"
)

func main() {
	cfg := config.Load()
	log := observability.NewLogger(cfg.Env)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("db connection failed", slog.Any("err", err))
		os.Exit(1)
	}
	defer pool.Close()

	repo := tenants.NewRepository(pool)
	svc := tenants.NewService(repo)
	server := apphttp.NewServer(log, svc)

	httpServer := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      server.Routes(cfg.ReadTimeout),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("http server started", slog.String("addr", cfg.HTTPAddr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", slog.Any("err", err))
			stop()
		}
	}()

	<-ctx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(shutdownCtx)

	log.Info("shutdown complete")
}
