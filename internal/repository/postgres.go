package repository

import (
	"context"
	"database/sql"
	"fmt"
	"l0/internal/config"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	db *sql.DB
}

func Connect(cfg config.Database) (*Storage, error) {
	const op = "storage.postgre.New"

	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", cfg.Host),
		getEnv("DB_PORT", cfg.Port),
		getEnv("DB_USER", cfg.User),
		getEnv("DB_PASSWORD", cfg.Password),
		getEnv("DB_NAME", cfg.Dbname),
	)

	log.Printf("Attempting to connect to database with")

	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return nil, fmt.Errorf("%s: goose up: %w", op, err)
	}

	migrationsDir, err := filepath.Abs("./migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get migrations path: %w", op, err)
	}

	log.Printf("Applying migrations from: %s", migrationsDir)
	if err := goose.Up(db, migrationsDir); err != nil {
		return nil, fmt.Errorf("%s: failed to apply migrations: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
