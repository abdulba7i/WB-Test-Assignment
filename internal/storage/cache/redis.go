package cache

import (
	"context"
	"l0wb/internal/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(cfg config.Redis) *Redis {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		Username: cfg.User,
		DB:       0, // use default DB
	})

	return &Redis{client: rdb}

}

func (r *Redis) Set(key string, value interface{}) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}
