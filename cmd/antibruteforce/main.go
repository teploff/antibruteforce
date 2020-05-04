package main

import (
	"context"
	"flag"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/endpoints/auth"
	"github.com/teploff/antibruteforce/infrastructure/logger"
	"github.com/teploff/antibruteforce/internal/implementaion/service"
	kitgrpc "github.com/teploff/antibruteforce/transport/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFile = flag.String("config", "./init/config_dev.yaml", "configuration file path")
	dev        = flag.Bool("dev", false, "dev mode")
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

	srv := service.NewAuthService()

	grpcServer := kitgrpc.NewGRPCServer(auth.MakeAuthEndpoints(srv),
		logger.NewZapSugarLogger(zapLogger, zapcore.ErrorLevel))
	go func() {
		if err = grpcServer.Serve(grpcListener); err != grpc.ErrServerStopped && err != nil {
			zapLogger.Fatal("gRPC serve error", zap.Error(err))
		}
	}()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer func() {
		// extra handling here
		cancel()
	}()
}
