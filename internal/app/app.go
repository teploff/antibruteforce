package app

import (
	"context"
	"errors"
	"net"
	nethttp "net/http"
	"time"

	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/internal/domain/repository"
	"github.com/teploff/antibruteforce/internal/endpoints/admin"
	"github.com/teploff/antibruteforce/internal/endpoints/auth"
	"github.com/teploff/antibruteforce/internal/implementation/middleware"
	"github.com/teploff/antibruteforce/internal/implementation/repository/bucket"
	"github.com/teploff/antibruteforce/internal/implementation/service"
	"github.com/teploff/antibruteforce/internal/infrastructure/logger"
	"github.com/teploff/antibruteforce/internal/limiter"
	kitgrpc "github.com/teploff/antibruteforce/internal/transport/grpc"
	"github.com/teploff/antibruteforce/internal/transport/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

const cancelTimeout = time.Millisecond * 500

// AppOption via application.
type Option func(*App)

// WithLogger adding logger option.
func WithLogger(l *zap.Logger) Option {
	return func(a *App) {
		a.logger = l
	}
}

// WithLeakyBuckets adding leaky buckets for logins, passwords and ips.
func WithLeakyBuckets(cfg config.RateLimiterConfig) Option {
	return func(a *App) {
		a.loginBucket = bucket.NewLeakyBucket(cfg.Login.Rate, cfg.Login.Interval, cfg.Login.ExpireTime)
		a.passwordBucket = bucket.NewLeakyBucket(cfg.Password.Rate, cfg.Password.Interval, cfg.Password.ExpireTime)
		a.ipBucket = bucket.NewLeakyBucket(cfg.IP.Rate, cfg.IP.Interval, cfg.IP.ExpireTime)
	}
}

// WithIPList adding ip list for admin panel.
func WithIPList(ipList repository.IPStorable) Option {
	return func(a *App) {
		a.ipList = ipList
	}
}

// App is application to encapsulate login to launch in main.
type App struct {
	cfg             config.Config
	loginBucket     repository.BucketStorable
	passwordBucket  repository.BucketStorable
	ipBucket        repository.BucketStorable
	ipList          repository.IPStorable
	logger          *zap.Logger
	stopCommandChan chan struct{}
}

// NewApp returns instance of app.
func NewApp(cfg config.Config, opts ...Option) *App {
	app := &App{
		cfg:             cfg,
		logger:          zap.NewNop(),
		stopCommandChan: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// Run lunch application.
func (a *App) Run() {
	gRPCListener, err := net.Listen("tcp", a.cfg.GRPCServer.Addr)
	if err != nil {
		a.logger.Fatal("gRPC listener", zap.Error(err))
	}

	rateLimiter := limiter.NewRateLimiter(a.loginBucket, a.passwordBucket, a.ipBucket, a.cfg.RateLimiter.GCTime)

	authSvc := service.NewAuthService(rateLimiter, a.ipList)
	adminSvc := middleware.NewPrometheusMiddleware(
		service.NewAdminService(a.ipList, a.loginBucket, a.passwordBucket, a.ipBucket))

	gRPCServer := kitgrpc.NewGRPCServer(auth.MakeAuthEndpoints(authSvc),
		logger.NewZapSugarLogger(a.logger, zapcore.ErrorLevel))

	router := http.NewHTTPServer(admin.MakeAdminEndpoints(adminSvc), a.logger)
	srv := &nethttp.Server{
		Addr:    a.cfg.HTTPServer.Addr,
		Handler: router,
	}

	go rateLimiter.RunGarbageCollector()

	go func() {
		if err = gRPCServer.Serve(gRPCListener); !errors.Is(err, grpc.ErrServerStopped) && err != nil {
			a.logger.Fatal("gRPC serve error", zap.Error(err))
		}
	}()

	go func() {
		if err = srv.ListenAndServe(); !errors.Is(err, nethttp.ErrServerClosed) && err != nil {
			a.logger.Fatal("http serve error", zap.Error(err))
		}
	}()

	<-a.stopCommandChan

	gRPCServer.GracefulStop()
	rateLimiter.Close()

	ctx, cancel := context.WithTimeout(context.Background(), cancelTimeout)
	if err = srv.Shutdown(ctx); err != nil {
		a.logger.Fatal("http shutdown error", zap.Error(err))
	}

	defer func() {
		// extra handling here
		cancel()
	}()
}

// Stop gracefully shutting down application.
func (a *App) Stop() {
	a.stopCommandChan <- struct{}{}
}
