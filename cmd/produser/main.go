package main

import (
	"encoding/json"
	"l0wb/internal/config"
	"l0wb/internal/storage/postgres"
	_nats "l0wb/pkg/nats"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

func main() {
	cfg := config.MustLoad()
	nc, err := _nats.New(cfg.NatsStreaming, "1")

	if err != nil {
		log.Fatal(err)
	}

	for {
		orderData, err := json.Marshal(GetRandomOrder())
		if err != nil {
			log.Fatal("failed to marshal order data:", err)
		}

		// Отправляем JSON через NATS
		err = nc.Publish("l0wb", orderData)
		if err != nil {
			log.Fatal("failed to publish message:", err)
		}

		log.Info("message sent")
	}

	nc.Close()
}

func GetRandomOrder() postgres.Order {
	id := uuid.NewString()
	order := postgres.Order{
		OrderUID:    id,
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: postgres.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: postgres.Payment{
			Transaction:  id,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []postgres.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NMID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SMID:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OOFShard:          "1",
	}

	return order
}
