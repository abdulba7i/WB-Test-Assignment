package main

import (
	"log/slog"

	"l0wb/internal/config"
	_nats "l0wb/pkg/nats"

	"github.com/labstack/gommon/log"
	"github.com/nats-io/stan.go"
)

func main() {
	cfg := config.MustLoad()
	nc, err := _nats.New(cfg.NatsStreaming, "2")

	if err != nil {
		log.Fatal(err)
	}

	nc.Consume("l0wb", func(msg *stan.Msg) {
		log.Info("got message", slog.String("data", string(msg.Data)))
		msg.Ack()
	})

	select {}
}
