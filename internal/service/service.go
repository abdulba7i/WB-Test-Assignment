package service

import (
	"l0/internal/cache"
	"l0/internal/model"
	"l0/internal/repository"
	"sync"

	"github.com/labstack/gommon/log"
)

type OrderService struct {
	Storage *repository.Storage
	Redis   *cache.Redis
}

func New(storage *repository.Storage, redis *cache.Redis) *OrderService {
	return &OrderService{
		Storage: storage,
		Redis:   redis,
	}
}

func (s *OrderService) GetOrderById(id string) (model.Order, error) {
	var order model.Order

	err := s.Redis.Get(id, &order)
	if err == nil {
		log.Info("got order from redis")
		return order, nil
	}

	log.Info("got order from db")

	order, err = s.Storage.GetOrderById(id)
	if err != nil {
		return order, err
	}

	err = s.Redis.Set(id, order)
	if err != nil {
		log.Warnf("failed to set order to redis: %v", err)
	}

	return order, nil
}

func (s *OrderService) LoadOrdersToCache() error {
	log.Info("load orders from db")

	const limit = 100
	offset := 0
	var wg sync.WaitGroup

	for {
		orders, err := s.Storage.GetAllOrders(limit, offset)
		if err != nil {
			log.Errorf("failed to load orders from db: %v", err)
			break
		}
		if len(orders) == 0 {
			break
		}

		log.Infof("loading %d orders to cache...", len(orders))

		for _, order := range orders {
			wg.Add(1)
			go func(order model.Order) {
				defer wg.Done()
				err := s.Redis.Set(order.OrderUID, order)
				if err != nil {
					log.Warnf("failed to cache order %s: %v", order.OrderUID, err)
				}
			}(order)
		}

		offset += limit
	}

	wg.Wait()
	log.Info("all orders loaded to cache")
	return nil
}
