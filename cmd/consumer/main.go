package main

import (
	"context"
	"encoding/json"
	"l0/internal/config"
	"l0/internal/storage/cache"
	"l0/internal/storage/postgres"
	_nats "l0/pkg/nats"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	// Загрузка конфигурации
	cfg := config.MustLoad()

	// Настройка логгера
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	log.Info("starting consumer")

	// Инициализация хранилищ
	storage, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}

	redisCache := cache.New(cfg.Redis)
	cacheService := cache.NewCacheService(redisCache, storage, cache.WithLogger(log))

	// Восстанавливаем кэш из базы данных
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	if err := cacheService.RestoreCache(ctx); err != nil {
		log.Error("failed to restore cache", slog.Any("error", err))
		cancel()
		os.Exit(1)
	}
	cancel()
	log.Info("cache restored successfully")

	// Подключение к NATS Streaming
	nc, err := _nats.New(cfg.NatsStreaming, "consumer")
	if err != nil {
		log.Error("failed to connect to nats", slog.Any("error", err))
		os.Exit(1)
	}
	defer nc.Close()

	// Канал для graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// WaitGroup для отслеживания обработки сообщений
	var wg sync.WaitGroup
	_, cancel = context.WithCancel(context.Background())

	// Подписка на сообщения
	sub, err := nc.Consume("l0", func(msg *stan.Msg) {
		wg.Add(1)
		defer wg.Done()

		log.Info("received message", slog.String("data", string(msg.Data)))

		// Преобразование msg.Data в структуру Order
		var order postgres.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Error("failed to unmarshal message",
				slog.Any("error", err),
				slog.String("data", string(msg.Data)),
			)
			return
		}

		// Сохранение в базу данных
		if err := storage.AddOrder(order); err != nil {
			log.Error("failed to save order to database",
				slog.Any("error", err),
				slog.String("order_id", order.OrderUID),
			)
			return
		}

		// Сохранение в Redis
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

		// Подтверждение обработки сообщения
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

	// Ожидание сигнала завершения
	<-shutdown
	log.Info("starting graceful shutdown")

	// Отмена контекста и ожидание завершения обработки сообщений
	cancel()

	// Отписка от NATS
	if err := sub.Unsubscribe(); err != nil {
		log.Error("failed to unsubscribe", slog.Any("error", err))
	}

	// Ожидание завершения обработки текущих сообщений с таймаутом
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
