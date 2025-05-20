package cache

import (
	"context"
	"encoding/json"
	"l0/internal/config"
	"log"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(cfg config.Redis) *Redis {

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Host + ":" + cfg.Port,
		DB:   0,
	})

	return &Redis{client: rdb}

}

func (r *Redis) Set(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Printf("Failed to marshal value: %v", err)
		return err
	}

	err = r.client.Set(context.Background(), key, jsonData, 0).Err()
	if err != nil {
		log.Printf("Failed to set key %s in Redis: %v", key, err)
	}
	return err
}

func (r *Redis) Get(key string, dest interface{}) error {
	jsonData, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("Failed to get key %s from Redis: %v", key, err)
		return err
	}

	err = json.Unmarshal([]byte(jsonData), dest)
	if err != nil {
		log.Printf("Failed to unmarshal value: %v", err)
	}
	return err
}
