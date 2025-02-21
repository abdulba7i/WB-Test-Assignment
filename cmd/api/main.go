package main

import (
	"l0wb/internal/config"
	"l0wb/internal/http-server/handlers/order"
	"l0wb/internal/storage/cache"
	"l0wb/internal/storage/postgres"
	"l0wb/internal/storage/services"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	redis := cache.New(cfg.Redis)
	storage, err := postgres.New(cfg.Database)
	order_service := services.New(*storage, *redis)
	if err != nil {
		logger.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}
	go order_service.LoadOrdersToCache()
	// if err != nil {
	// 	logger.Error("failed load orders to cache", slog.Any("error", err))
	// 	os.Exit(1)
	// }

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/order", func(r chi.Router) {
		r.Get("/{id}", order.GetOrder(log, order_service))
	})

	log.Info("starting server", slog.String("addres", cfg.HTTPServer.Address))
	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(1)
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		// log = setupPrettySlog() // здесь преукрасили вывод логов для удобства
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	default: // моё дополнение
		panic("not supported env")
	}

	return log
}

// func setupPrettySlog() *slog.Logger {
// 	opts := slogpretty.PrettyHandlerOptions{
// 		SlogOpts: &slog.HandlerOptions{
// 			Level: slog.LevelDebug,
// 		},
// 	}

// 	handler := opts.NewPrettyHandler(os.Stdout)

// 	return slog.New(handler)
// }
