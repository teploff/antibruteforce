package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/endpoints/auth"
	"github.com/teploff/antibruteforce/infrastructure/logger"
	"github.com/teploff/antibruteforce/internal/implementation/gateway"
	"github.com/teploff/antibruteforce/internal/implementation/service"
	kitgrpc "github.com/teploff/antibruteforce/transport/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

const timeoutDuration = 500 * time.Millisecond

var (
	//nolint:gochecknoglobals
	configFile = flag.String("config", "./init/config_dev.yaml", "configuration file path")
	//nolint:gochecknoglobals
	dev = flag.Bool("dev", false, "dev mode")
)

func main() {
	flag.Parse()

	cfg, err := config.LoadFromFile(*configFile)
	if err != nil {
		panic(err)
	}

	zapLogger := logger.New(*dev, &cfg.Logger)

	grpcListener, err := net.Listen("tcp", cfg.Server.Addr)
	if err != nil {
		zapLogger.Fatal("gRPC listener", zap.Error(err))
	}

	srv := service.NewAuthService(gateway.NewRateLimiter(cfg.Limiter.Login),
		gateway.NewRateLimiter(cfg.Limiter.Password), gateway.NewRateLimiter(cfg.Limiter.IP))

	grpcServer := kitgrpc.NewGRPCServer(auth.MakeAuthEndpoints(srv),
		logger.NewZapSugarLogger(zapLogger, zapcore.ErrorLevel))

	go func() {
		if err = grpcServer.Serve(grpcListener); errors.Is(err, grpc.ErrServerStopped) && err != nil {
			zapLogger.Fatal("gRPC serve error", zap.Error(err))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	_, cancel := context.WithTimeout(context.Background(), timeoutDuration)

	defer func() {
		// extra handling here
		cancel()
	}()
}
