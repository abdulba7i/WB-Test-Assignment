package services

import (
	"l0wb/internal/storage/cache"
	"l0wb/internal/storage/postgres"

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
