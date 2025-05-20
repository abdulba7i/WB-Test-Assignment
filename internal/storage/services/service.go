package services

import (
	"l0/internal/storage/cache"
	"l0/internal/storage/postgres"

	// "l0/internal/repository/cache"
	// "l0/internal/postgrte/postgres"

	"github.com/labstack/gommon/log"
)

type OrderService struct {
	Storage *postgres.Storage
	Redis   *cache.Redis
}

func New(storage postgres.Storage, redis cache.Redis) *OrderService {
	return &OrderService{Storage: &storage, Redis: &redis}
}

func (s *OrderService) GetOrderById(id string) (postgres.Order, error) {
	var order postgres.Order
	err := s.Redis.Get(id, &order)
	if err != nil {
		log.Info("got order from db")
		order, err = s.Storage.GetOrderById(id)
		if err != nil {
			return order, err
		}
		err = s.Redis.Set(id, order)
		if err != nil {
			return order, err
		}
		return order, nil
	}
	log.Info("got order from redis")
	return order, nil
}

func (s *OrderService) LoadOrdersToCache() error {
	log.Info("load orders from db")
	limit := 100
	ofset := 0
	for {
		orders, err := s.Storage.GetAllOrders(limit, ofset)
		log.Info("got orders from db")
		if err != nil {
			log.Error("failed load orders to cache", err)
			continue
		}
		if len(orders) == 0 {
			break
		}
		log.Info("loading orders to cache")
		for _, order := range orders {
			go s.Redis.Set(order.OrderUID, order)
			// if err != nil {
			// 	continue
			// }
		}
		limit += 100
		ofset += 100
	}

	return nil
}
