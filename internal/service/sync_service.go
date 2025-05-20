package service

import (
	"context"
	"fmt"

	// "l0/internal/storage/cache"

	// "l0/garbage/storage/postgres"
	"l0/internal/cache"
	"l0/internal/model"
	"l0/internal/repository"

	// "l0/internal/storage/postgres"
	"log"
	"time"
)

type SyncService struct {
	db          *repository.Storage
	cache       *cache.CacheService
	batchSize   int
	syncTimeout time.Duration
}

func NewSyncService(db *repository.Storage, cache *cache.CacheService, batchSize int, syncTimeout time.Duration) *SyncService {
	return &SyncService{
		db:          db,
		cache:       cache,
		batchSize:   batchSize,
		syncTimeout: syncTimeout,
	}
}

// SyncData синхронизирует данные между PostgreSQL и Redis
func (s *SyncService) SyncData(ctx context.Context) error {
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, s.syncTimeout)
	defer cancel()

	// Очищаем старые данные из кэша
	if err := s.cache.ClearOldData(ctx, "*"); err != nil {
		log.Printf("Failed to clear old data: %v", err)
	}

	// Функция для обработки пакета заказов
	processBatch := func(orders []model.Order) error {
		// Преобразуем заказы в интерфейс для кэша
		var data []interface{}
		for _, order := range orders {
			data = append(data, order)
		}

		// Загружаем данные в кэш
		return s.cache.LoadDataBatch(ctx, data, func(item interface{}) string {
			order := item.(model.Order)
			return fmt.Sprintf("order:%s", order.OrderUID)
		})
	}

	// Получаем и обрабатываем данные пакетами
	return s.db.GetOrdersBatch(ctx, s.batchSize, processBatch)
}

// StartPeriodicSync запускает периодическую синхронизацию
func (s *SyncService) StartPeriodicSync(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.SyncData(ctx); err != nil {
				log.Printf("Failed to sync data: %v", err)
			}
		}
	}
}
