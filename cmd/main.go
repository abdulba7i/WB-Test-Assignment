package main

import (
	"l0wb/internal/config"
)

func main() {
	/* cfg */ _ = config.MustLoad()
}

// func setupLogger(env string) *slog.Logger {
// 	var log *slog.Logger

// }
