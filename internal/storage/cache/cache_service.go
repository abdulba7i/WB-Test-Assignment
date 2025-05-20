package cache

import (
	"context"
	"fmt"
	"l0/internal/storage/postgres"
	"log"
	"log/slog"
	"sync"
	"time"
)

const (
	defaultBatchSize = 1000
	defaultWorkers   = 5
)

type CacheService struct {
	redis      *Redis
	storage    *postgres.Storage
	logger     *slog.Logger
	batchSize  int
	workers    int
	expiration time.Duration
}

type CacheServiceOption func(*CacheService)

func WithBatchSize(size int) CacheServiceOption {
	return func(cs *CacheService) {
		if size > 0 {
			cs.batchSize = size
		}
	}
}

func WithWorkers(workers int) CacheServiceOption {
	return func(cs *CacheService) {
		if workers > 0 {
			cs.workers = workers
		}
	}
}

func WithExpiration(exp time.Duration) CacheServiceOption {
	return func(cs *CacheService) {
		cs.expiration = exp
	}
}

func WithLogger(logger *slog.Logger) CacheServiceOption {
	return func(cs *CacheService) {
		cs.logger = logger
	}
}

func NewCacheService(redis *Redis, storage *postgres.Storage, opts ...CacheServiceOption) *CacheService {
	cs := &CacheService{
		redis:      redis,
		storage:    storage,
		logger:     slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{Level: slog.LevelInfo})),
		batchSize:  defaultBatchSize,
		workers:    defaultWorkers,
		expiration: 24 * time.Hour, // по умолчанию 24 часа
	}

	for _, opt := range opts {
		opt(cs)
	}

	return cs
}

// RestoreCache восстанавливает кэш из базы данных
func (cs *CacheService) RestoreCache(ctx context.Context) error {
	cs.logger.Info("starting cache restoration from database")

	// Получаем все заказы из базы данных
	orders, err := cs.storage.GetAllOrders(1000, 0) // Получаем первую тысячу заказов
	if err != nil {
		return fmt.Errorf("failed to get orders from database: %w", err)
	}

	cs.logger.Info("retrieved orders from database", slog.Int("count", len(orders)))

	// Создаем каналы для распределения работы
	jobs := make(chan postgres.Order, len(orders))
	errors := make(chan error, len(orders))
	var wg sync.WaitGroup

	// Запускаем воркеры
	for i := 0; i < cs.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for order := range jobs {
				if err := cs.redis.Set(order.OrderUID, order); err != nil {
					errors <- fmt.Errorf("failed to cache order %s: %w", order.OrderUID, err)
					continue
				}
			}
		}()
	}

	// Отправляем заказы воркерам
	for _, order := range orders {
		select {
		case jobs <- order:
		case <-ctx.Done():
			close(jobs)
			return ctx.Err()
		}
	}

	// Закрываем канал работ и ждем завершения всех воркеров
	close(jobs)
	wg.Wait()

	// Проверяем ошибки
	select {
	case err := <-errors:
		return fmt.Errorf("error during cache restoration: %w", err)
	default:
		cs.logger.Info("cache restoration completed successfully", slog.Int("orders_cached", len(orders)))
		return nil
	}
}

// LoadDataBatch загружает данные в Redis пакетами
func (cs *CacheService) LoadDataBatch(ctx context.Context, data []interface{}, keyFunc func(interface{}) string) error {
	if len(data) == 0 {
		return nil
	}

	// Создаем каналы для распределения работы
	jobs := make(chan interface{}, len(data))
	errors := make(chan error, len(data))
	var wg sync.WaitGroup

	// Запускаем воркеры
	for i := 0; i < cs.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				key := keyFunc(item)
				if err := cs.redis.Set(key, item); err != nil {
					errors <- fmt.Errorf("failed to cache item %s: %w", key, err)
					continue
				}
			}
		}()
	}

	// Отправляем данные воркерам
	for _, item := range data {
		select {
		case jobs <- item:
		case <-ctx.Done():
			close(jobs)
			return ctx.Err()
		}
	}

	// Закрываем канал работ и ждем завершения всех воркеров
	close(jobs)
	wg.Wait()

	// Проверяем ошибки
	select {
	case err := <-errors:
		return err
	default:
		return nil
	}
}

// GetCachedData получает данные из кэша
func (cs *CacheService) GetCachedData(key string, dest interface{}) error {
	return cs.redis.Get(key, dest)
}

// ClearOldData очищает старые данные из кэша
func (cs *CacheService) ClearOldData(ctx context.Context, pattern string) error {
	iter := cs.redis.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if err := cs.redis.client.Del(ctx, key).Err(); err != nil {
			cs.logger.Error("failed to delete key",
				slog.String("key", key),
				slog.Any("error", err),
			)
		}
	}
	return iter.Err()
}
