package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/qosmioo/ulr-shortener/config"
	"github.com/qosmioo/ulr-shortener/internal/delivery"
	"github.com/qosmioo/ulr-shortener/internal/storage"
	"github.com/qosmioo/ulr-shortener/internal/usecase"
	pb "github.com/qosmioo/ulr-shortener/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := initLogger()
	defer logger.Sync()

	cfg := loadConfig(logger)

	store := initStorage(cfg, logger)

	urlShortenerService := usecase.NewURLShortenerService(store)

	grpcServer := startGRPCServer(cfg, urlShortenerService, logger)
	httpSrv := startHTTPServer(cfg, urlShortenerService, logger)

	gracefulShutdown(grpcServer, httpSrv, logger)
}

func initLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func loadConfig(logger *zap.Logger) *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}
	return cfg
}

func initStorage(cfg *config.Config, logger *zap.Logger) storage.Storage {
	var store storage.Storage
	var err error

	logger.Info("Initializing storage", zap.String("type", cfg.Database.Type))

	switch cfg.Database.Type {
	case "postgres":
		store, err = storage.NewPostgresStorage(cfg.Database.Postgres, logger)
		if err != nil {
			logger.Fatal("Failed to initialize Postgres storage", zap.Error(err))
		}
	case "redis":
		store, err = storage.NewRedisStorage(cfg.Database.Redis, logger)
		if err != nil {
			logger.Fatal("Failed to initialize Redis storage", zap.Error(err))
		}
	default:
		logger.Fatal("Unsupported storage type")
	}

	return store
}

func startGRPCServer(cfg *config.Config, usecase *usecase.URLShortenerService, logger *zap.Logger) *grpc.Server {
	grpcLis, err := net.Listen("tcp", ":"+cfg.Server.GrpcPort)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGRPCHandlerServer(grpcServer, delivery.NewGRPCServer(usecase, logger))
	reflection.Register(grpcServer)

	go func() {
		logger.Info("gRPC server listening", zap.String("address", grpcLis.Addr().String()))
		if err := grpcServer.Serve(grpcLis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	return grpcServer
}

func startHTTPServer(cfg *config.Config, usecase *usecase.URLShortenerService, logger *zap.Logger) *http.Server {
	router := mux.NewRouter()
	httpServer := delivery.NewHTTPServer(usecase, logger)
	httpServer.ApiEndpoints(router)

	httpSrv := &http.Server{
		Addr:    ":" + cfg.Server.HttpPort,
		Handler: router,
	}

	go func() {
		logger.Info("HTTP server listening", zap.String("port", cfg.Server.HttpPort))
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to serve HTTP", zap.Error(err))
		}
	}()

	return httpSrv
}

func gracefulShutdown(grpcServer *grpc.Server, httpSrv *http.Server, logger *zap.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	logger.Info("Shutting down servers...")

	grpcServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		logger.Fatal("HTTP server shutdown failed", zap.Error(err))
	}

	logger.Info("Servers stopped gracefully")
}
