package main

import (
	"fmt"
	"l0wb/internal/config"
	_nats "l0wb/pkg/nats"

	"github.com/labstack/gommon/log"
)

func main() {
	cfg := config.MustLoad()
	nc, err := _nats.New(cfg.NatsStreaming, "1")

	if err != nil {
		log.Fatal(err)
	}

	err = nc.Publish("l0wb", []byte(fmt.Sprintf("test message %d", 1)))
	err = nc.Publish("l0wb", []byte(fmt.Sprintf("test message %d", 2)))

	// for i := 0; i < 2; i++ {
	// 	err = nc.Publish("l0wb", []byte(fmt.Sprintf("test message %d", i)))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	log.Info("message sent")
	// }

	nc.Close()
}
