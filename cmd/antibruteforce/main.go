package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/infrastructure/logger"
	"github.com/teploff/antibruteforce/internal/implementation/repository/ip"
	"github.com/teploff/antibruteforce/pkg"
	"go.uber.org/zap"
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

	ipList, err := ip.NewRedisIPList(cfg.Redis)
	if err != nil {
		zapLogger.Fatal("redis connect error", zap.Error(err))
	}

	app := pkg.NewApp(cfg,
		pkg.WithLogger(zapLogger),
		pkg.WithLeakyBuckets(cfg.RateLimiter),
		pkg.WithIPList(ipList))

	app.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	app.Stop()
}
