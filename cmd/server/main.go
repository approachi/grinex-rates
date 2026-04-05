package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/approachi/grinex-rates/internal/config"
	"github.com/approachi/grinex-rates/internal/exchange"
	grpcserver "github.com/approachi/grinex-rates/internal/grpc"
	"github.com/approachi/grinex-rates/internal/observability/logger"
	"github.com/approachi/grinex-rates/internal/observability/metrics"
	"github.com/approachi/grinex-rates/internal/observability/tracing"
	"github.com/approachi/grinex-rates/internal/repository/postgres"
	"github.com/approachi/grinex-rates/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapLogger, err := logger.New(cfg.Log.Level)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer zapLogger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, err := tracing.NewTracerProvider(ctx, cfg.Tracing.Endpoint)
	if err != nil {
		zapLogger.Warn("failed to init tracer, tracing disabled", zap.Error(err))
	} else {
		defer tp.Shutdown(ctx)
	}

	metrics.Register()
	go func() {
		zapLogger.Info("starting metrics server", zap.Int("port", cfg.Metrics.Port))
		if err := metrics.RunServer(ctx, cfg.Metrics.Port); err != nil {
			zapLogger.Warn("metrics server stopped", zap.Error(err))
		}
	}()

	pool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		zapLogger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer pool.Close()

	repo := postgres.NewRepository(pool)
	client := exchange.NewClient()
	svc := service.NewRatesService(client, repo)
	handler := grpcserver.NewHandler(svc)
	srv := grpcserver.NewServer(handler, cfg.GRPC.Port)

	go func() {
		zapLogger.Info("starting gRPC server", zap.Int("port", cfg.GRPC.Port))
		if err := srv.Run(); err != nil {
			zapLogger.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("shutting down gRPC server")
	srv.Stop()
	zapLogger.Info("server stopped")
}
