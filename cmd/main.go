package main

import (
	"config_center/internal/application"
	"config_center/internal/config"
	"context"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("an error occurred during intialization of the logger: %v", err)
	}

	defer func() {
		if err = logger.Sync(); err != nil {
			log.Fatalf("an error occurred during flushing the logs buffer: %v", err)
		}
	}()

	logger.Info("config_center start")

	cfg, err := config.FromEnv()
	if err != nil {
		//log.Fatalf("a fatal error occurred during reading env config: %v", err)
		logger.Error("a fatal error occurred during reading env config", zap.Error(err))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	termSignalChan := make(chan os.Signal, 1)
	signal.Notify(termSignalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	app := application.New(cfg, logger)
	errChan := make(chan error)
	wg, gracefulQuitFn := application.NewAutoWaitGroup()

	err = app.Start(ctx, errChan, gracefulQuitFn)
	if err != nil {
		log.Fatalf("an error occurred during start of the application: %v", err)
	}

	select {
	case sig := <-termSignalChan:
		logger.Info("received system signal", zap.String("signal", sig.String()))
		cancel()
	case fatalErr := <-errChan:
		logger.Error("a fatal error occurred during work of the service", zap.Error(fatalErr))
		cancel()
	}

	wg.Wait()

	logger.Info("config_center exit")
}
