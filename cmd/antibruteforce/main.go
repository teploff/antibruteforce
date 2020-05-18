package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/internal/app"
	"github.com/teploff/antibruteforce/internal/implementation/repository/ip"
	"github.com/teploff/antibruteforce/internal/infrastructure/logger"
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

	ipList, err := ip.NewMongoIPList(cfg.Mongo)
	if err != nil {
		zapLogger.Fatal("mongodb connect error", zap.Error(err))
	}

	application := app.NewApp(cfg,
		app.WithLogger(zapLogger),
		app.WithLeakyBuckets(cfg.RateLimiter),
		app.WithIPList(ipList))

	application.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
