package main

import (
	"l0wb/internal/config"
	"l0wb/internal/storage/postgres"
	"log/slog"
	"os"

	"github.com/labstack/gommon/log"
)

func main() {
	cfg := config.MustLoad()

	log.Info(
		"starting url-shortener", slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")
	_, err := postgres.New(cfg.Database)

	if err != nil {
		log.Error("failed to init storage", "error: %v", err)
		os.Exit(1)
	}
}
