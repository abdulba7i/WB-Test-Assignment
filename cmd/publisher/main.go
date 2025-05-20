package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"l0/internal/config"
	"log/slog"
	"os"
	"time"

	"github.com/nats-io/stan.go"
)

var testOrder = `{
	"order_uid": "b563feb7b2b84best",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
		"name": "Test Tev",
		"phone": "+9720000000",
		"zip": "2639809",
		"city": "Kiryat Mozkin",
		"address": "Ploshad Mira 15",
		"region": "Kraiot",
		"email": "test@gmail.com"
	},
	"payment": {
		"transaction": "b563feb7b2b84b6test",
		"request_id": "",
		"currency": "USD",
		"provider": "wbpay",
		"amount": 1817,
		"payment_dt": 1637907727,
		"bank": "alpha",
		"delivery_cost": 1500,
		"goods_total": 317,
		"custom_fee": 0
	},
	"items": [
		{
			"chrt_id": 9934930,
			"track_number": "WBILMTESTTRACK",
			"price": 453,
			"rid": "ab4219087a764ae0btest",
			"name": "Mascaras",
			"sale": 30,
			"size": "0",
			"total_price": 317,
			"nm_id": 2389212,
			"brand": "Vivienne Sabo",
			"status": 202
		}
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
}`

func main() {
	// Флаги командной строки
	count := flag.Int("n", 10, "количество сообщений для отправки")
	delay := flag.Duration("delay", 1*time.Second, "задержка между сообщениями")
	flag.Parse()

	// Загрузка конфигурации
	cfg := config.MustLoad()

	// Настройка логгера
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Формируем URL для подключения к NATS
	natsURL := fmt.Sprintf("nats://%s:%s", cfg.NatsStreaming.Host, cfg.NatsStreaming.Port)
	if cfg.NatsStreaming.User != "" && cfg.NatsStreaming.Password != "" {
		natsURL = fmt.Sprintf("nats://%s:%s@%s:%s",
			cfg.NatsStreaming.User,
			cfg.NatsStreaming.Password,
			cfg.NatsStreaming.Host,
			cfg.NatsStreaming.Port,
		)
	}

	// Подключение к NATS Streaming
	sc, err := stan.Connect(
		cfg.NatsStreaming.ClusterID,
		"publisher",
		stan.NatsURL(natsURL),
	)
	if err != nil {
		log.Error("failed to connect to nats", slog.Any("error", err))
		os.Exit(1)
	}
	defer sc.Close()

	// Проверка валидности JSON
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(testOrder), &jsonMap); err != nil {
		log.Error("invalid test order JSON", slog.Any("error", err))
		os.Exit(1)
	}

	// Отправка сообщений
	// for i := 0; i < *count; i++ {
	// 	jsonMap["order_uid"] = fmt.Sprintf("b563feb7b2b84best%d", time.Now().Unix())
	for i := 0; i < *count; i++ {
		jsonMap["order_uid"] = fmt.Sprintf("b563feb7b2b84best-%d-%d", time.Now().UnixMilli(), i)
		jsonData, err := json.Marshal(jsonMap)
		if err != nil {
			log.Error("failed to marshal test order", slog.Any("error", err))
			os.Exit(1)
		}
		if err := sc.Publish("l0", jsonData); err != nil {
			log.Error("failed to publish message",
				slog.Any("error", err),
				slog.Int("message_number", i+1),
			)
			continue
		}

		log.Info("message published successfully",
			slog.Int("message_number", i+1),
			slog.Int("total_messages", *count),
		)

		if i < *count-1 { // Не ждем после последнего сообщения
			time.Sleep(*delay)
		}
	}
}
