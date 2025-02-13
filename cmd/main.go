package main

import (
	"encoding/json"
	"l0wb/internal/config"
	"l0wb/internal/storage/postgres"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)

	storage, err := postgres.New(cfg.Database)
	if err != nil {
		logger.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}

	result, err := storage.GetOrderById("b563feb7b2b84b6test")
	if err != nil {
		logger.Error("failed to get order", slog.Any("error", err))
		os.Exit(1)
	}

	jsonOrder, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logger.Error("failed to serialize order", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("retrieved order", slog.String("order", string(jsonOrder)))
}
