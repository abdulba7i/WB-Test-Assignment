package service

import (
	"context"
	"fmt"
	"l0/internal/cache"
	"l0/internal/model"
	"l0/internal/repository"

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

func (s *SyncService) SyncData(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.syncTimeout)
	defer cancel()

	if err := s.cache.ClearOldData(ctx, "*"); err != nil {
		log.Printf("Failed to clear old data: %v", err)
	}

	processBatch := func(orders []model.Order) error {
		var data []interface{}
		for _, order := range orders {
			data = append(data, order)
		}

		return s.cache.LoadDataBatch(ctx, data, func(item interface{}) string {
			order := item.(model.Order)
			return fmt.Sprintf("order:%s", order.OrderUID)
		})
	}

	return s.db.GetOrdersBatch(ctx, s.batchSize, processBatch)
}

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
