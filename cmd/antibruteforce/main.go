package main

import (
	"context"
	"errors"
	"flag"
	"net"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/endpoints/admin"
	"github.com/teploff/antibruteforce/endpoints/auth"
	"github.com/teploff/antibruteforce/infrastructure/logger"
	"github.com/teploff/antibruteforce/internal/implementation/repository/bucket"
	"github.com/teploff/antibruteforce/internal/implementation/repository/ip"
	"github.com/teploff/antibruteforce/internal/implementation/service"
	"github.com/teploff/antibruteforce/internal/limiter"
	kitgrpc "github.com/teploff/antibruteforce/transport/grpc"
	"github.com/teploff/antibruteforce/transport/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

const cancelTimeout = time.Millisecond * 500

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

	gRPCListener, err := net.Listen("tcp", cfg.GRPCServer.Addr)
	if err != nil {
		zapLogger.Fatal("gRPC listener", zap.Error(err))
	}

	loginBuckets := bucket.NewLeakyBucket(cfg.RateLimiter.Login.Rate, cfg.RateLimiter.Login.Interval,
		cfg.RateLimiter.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfg.RateLimiter.Password.Rate, cfg.RateLimiter.Password.Interval,
		cfg.RateLimiter.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfg.RateLimiter.IP.Rate, cfg.RateLimiter.IP.Interval,
		cfg.RateLimiter.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfg.RateLimiter.GCTime)
	ipList, err := ip.NewRedisIPList(cfg.Redis)

	if err != nil {
		zapLogger.Fatal("redis connect error", zap.Error(err))
	}

	go rateLimiter.RunGarbageCollector()

	authSvc := service.NewAuthService(rateLimiter, ipList)
	adminSvc := service.NewAdminService(ipList, loginBuckets, passwordBuckets, ipBuckets)

	gRPCServer := kitgrpc.NewGRPCServer(auth.MakeAuthEndpoints(authSvc),
		logger.NewZapSugarLogger(zapLogger, zapcore.ErrorLevel))

	router := http.NewHTTPServer(admin.MakeAdminEndpoints(adminSvc), zapLogger)
	srv := &nethttp.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	go func() {
		if err = gRPCServer.Serve(gRPCListener); !errors.Is(err, grpc.ErrServerStopped) && err != nil {
			zapLogger.Fatal("gRPC serve error", zap.Error(err))
		}
	}()

	go func() {
		if err = srv.ListenAndServe(); !errors.Is(err, nethttp.ErrServerClosed) && err != nil {
			zapLogger.Fatal("http serve error", zap.Error(err))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	gRPCServer.GracefulStop()
	rateLimiter.Close()

	ctx, cancel := context.WithTimeout(context.Background(), cancelTimeout)
	if err = srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("http shutdown error", zap.Error(err))
	}

	defer func() {
		// extra handling here
		cancel()
		time.Sleep(time.Second)
	}()
}
