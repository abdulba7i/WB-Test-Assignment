package main

import (
	"context"
	"encoding/json"

	"l0/internal/cache"
	"l0/internal/config"
	"l0/internal/model"
	_nats "l0/internal/pkg/nats"
	"l0/internal/repository"

	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	log.Info("starting consumer")

	storage, err := repository.Connect(cfg.Database)
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}

	redisCache := cache.New(cfg.Redis)
	cacheService := cache.NewCacheService(redisCache, storage, cache.WithLogger(log))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	if err := cacheService.RestoreCache(ctx); err != nil {
		log.Error("failed to restore cache", slog.Any("error", err))
		cancel()
		os.Exit(1)
	}
	cancel()
	log.Info("cache restored successfully")

	nc, err := _nats.New(cfg.NatsStreaming, "consumer")
	if err != nil {
		log.Error("failed to connect to nats", slog.Any("error", err))
		os.Exit(1)
	}
	defer nc.Close()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	_, cancel = context.WithCancel(context.Background())

	sub, err := nc.Consume("l0", func(msg *stan.Msg) {
		wg.Add(1)
		defer wg.Done()

		log.Info("received message", slog.String("data", string(msg.Data)))

		var order model.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Error("failed to unmarshal message",
				slog.Any("error", err),
				slog.String("data", string(msg.Data)),
			)
			return
		}

		if err := storage.AddOrder(order); err != nil {
			log.Error("failed to save order to database",
				slog.Any("error", err),
				slog.String("order_id", order.OrderUID),
			)
			return
		}

		if err := redisCache.Set(order.OrderUID, order); err != nil {
			log.Error("failed to save order to redis",
				slog.Any("error", err),
				slog.String("order_id", order.OrderUID),
			)
			return
		}

		log.Info("order processed successfully",
			slog.String("order_id", order.OrderUID),
		)

		if err := msg.Ack(); err != nil {
			log.Error("failed to acknowledge message",
				slog.Any("error", err),
				slog.String("order_id", order.OrderUID),
			)
		}
	}, stan.DurableName("my-durable"), stan.SetManualAckMode())

	if err != nil {
		log.Error("failed to subscribe", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("consumer started successfully")

	<-shutdown
	log.Info("starting graceful shutdown")

	cancel()

	if err := sub.Unsubscribe(); err != nil {
		log.Error("failed to unsubscribe", slog.Any("error", err))
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info("all messages processed")
	case <-time.After(30 * time.Second):
		log.Warn("shutdown timeout exceeded")
	}

	log.Info("consumer stopped")
}
