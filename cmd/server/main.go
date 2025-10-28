package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/ionut-maxim/insider-project/internal/server"
	"github.com/ionut-maxim/insider-project/internal/store/postgres"
	"github.com/ionut-maxim/insider-project/internal/worker"
	"github.com/ionut-maxim/insider-project/internal/worker/cache/redis"
	"github.com/ionut-maxim/insider-project/internal/worker/notifier/webhook"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	var cfg config
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: "APP_"}); err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.LogLevel}))

	store, err := postgres.New(cfg.Database.ConnectionString(), logger)
	if err != nil {
		logger.Error("Failed to connect to database", "details", err)
		os.Exit(1)
	}

	notifier := webhook.New(cfg.WebhookURL)

	var cacheOptions []redis.Options
	if cfg.Cache.TTL != time.Duration(0) {
		cacheOptions = append(cacheOptions, redis.WithTTL(cfg.Cache.TTL))
	}
	cache := redis.New(cfg.Cache.ConnectionString(), cfg.Cache.Database, cacheOptions...)

	var workerOptions []worker.Option
	if cfg.PollInterval != time.Duration(0) {
		workerOptions = append(workerOptions, worker.WithInterval(cfg.PollInterval))
	}
	w := worker.New(store, notifier, cache, logger, workerOptions...)

	if err = w.Start(ctx); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	srv := server.New(cfg.ServerPort, store, w, logger)
	if err = srv.Start(); err != nil {
		logger.Error("Failed to start server", "details", err)
		os.Exit(1)
	}

	// Wait for CTRL-C
	<-ctx.Done()

	if err = srv.Close(); err != nil {
		logger.Error("Error closing server", "details", err)
		os.Exit(1)
	}
}
