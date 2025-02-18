package main

import (
	"encoding/json"
	"log/slog"

	"l0wb/internal/config"
	"l0wb/internal/storage/cache"
	"l0wb/internal/storage/postgres"
	_nats "l0wb/pkg/nats"

	"github.com/labstack/gommon/log"
	"github.com/nats-io/stan.go"
)

func main() {
	cfg := config.MustLoad()
	nc, err := _nats.New(cfg.NatsStreaming, "2")

	storage, err := postgres.New(cfg.Database)
	redis := cache.New(cfg.Redis)

	if err != nil {
		log.Fatal(err)
	}

	nc.Consume("l0wb", func(msg *stan.Msg) {
		log.Info("got message", slog.String("data", string(msg.Data)))

		// Преобразование msg.Data в структуру Order
		var order postgres.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Error("failed to unmarshal msg data:", err)
		}

		// Вызываем метод добавления заказа в БД
		err = storage.AddOrder(order)
		if err != nil {
			log.Error("failed add order	to DB:", err)
		}
		err = redis.Set(order.OrderUID, order)
		if err != nil {
			log.Error("failed add order to Redis:", err)
		}

		// Подтверждаем получение сообщения
		msg.Ack()
	})

	select {}
}
