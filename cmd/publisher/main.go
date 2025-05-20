package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"l0/internal/config"
	"l0/internal/lib/utils"
	"log/slog"
	"os"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	count := flag.Int("n", 10, "количество сообщений для отправки")
	delay := flag.Duration("delay", 1*time.Second, "задержка между сообщениями")
	flag.Parse()

	cfg := config.MustLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	natsURL := fmt.Sprintf("nats://%s:%s", cfg.NatsStreaming.Host, cfg.NatsStreaming.Port)
	if cfg.NatsStreaming.User != "" && cfg.NatsStreaming.Password != "" {
		natsURL = fmt.Sprintf("nats://%s:%s@%s:%s",
			cfg.NatsStreaming.User,
			cfg.NatsStreaming.Password,
			cfg.NatsStreaming.Host,
			cfg.NatsStreaming.Port,
		)
	}

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

	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(utils.TestOrder), &jsonMap); err != nil {
		log.Error("invalid test order JSON", slog.Any("error", err))
		os.Exit(1)
	}

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

		if i < *count-1 {
			time.Sleep(*delay)
		}
	}
}
