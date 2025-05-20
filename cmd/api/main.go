package main

import (
	"context"
	"l0/internal/config"
	"l0/internal/http-server/handlers/order"
	"l0/internal/storage/cache"
	"l0/internal/storage/postgres"
	"l0/internal/storage/services"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Загрузка конфигурации
	cfg := config.MustLoad()

	// Настройка логгера
	log := setupLogger(cfg.Env)
	log.Info("starting api server", slog.String("env", cfg.Env))

	// Инициализация хранилищ
	redis := cache.New(cfg.Redis)
	storage, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}

	// Инициализация сервиса
	orderService := services.New(*storage, *redis)

	// Настройка роутера
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Регистрация маршрутов
	router.Route("/order", func(r chi.Router) {
		r.Get("/{id}", order.GetOrder(log, orderService))
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(router)

	// Создание HTTP сервера
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      handler,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Канал для graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в горутине
	serverError := make(chan error, 1)
	go func() {
		log.Info("starting http server", slog.String("address", cfg.HTTPServer.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
	}()

	// Ожидание сигнала завершения или ошибки
	select {
	case err := <-serverError:
		log.Error("server error", slog.Any("error", err))
	case sig := <-shutdown:
		log.Info("starting shutdown", slog.String("signal", sig.String()))

		// Создание контекста с таймаутом для graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("failed to stop server", slog.Any("error", err))
			if err := srv.Close(); err != nil {
				log.Error("failed to close server", slog.Any("error", err))
			}
		}
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("unsupported environment: " + env)
	}
	return log
}
