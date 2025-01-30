package postgres

func Save()

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"time"
// 	"url-shortener/internal/storage"

// 	"github.com/jackc/pgx/v5/pgconn"
// 	_ "github.com/lib/pq"
// )

// type Storage struct {
// 	db *sql.DB
// }

// type DB struct {
// 	Host     string `yaml:"host"`
// 	Port     string `yaml:"port"`
// 	User     string `yaml:"user"`
// 	Password string `yaml:"password"`
// 	Dbname   string `yaml:"dbname"`
// }

// func New(c DB) (*Storage, error) {
// 	const op = "storage.postgre.New"

// 	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Dbname)

// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	// Создание таблицы
// 	stmt1 := `
// 	CREATE TABLE IF NOT EXISTS url (
// 		id SERIAL PRIMARY KEY,
// 		alias TEXT NOT NULL UNIQUE,
// 		url TEXT NOT NULL
// 	);`
// 	_, err = db.Exec(stmt1)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	// Создание индекса
// 	stmt2 := `
// 	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`
// 	_, err = db.Exec(stmt2)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return &Storage{db: db}, nil
// }

// func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
// 	const op = "storage.postgre.SaveUrl"

// 	var id int64
// 	query := "INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id"

// 	err := s.db.QueryRow(query, urlToSave, alias).Scan(&id)
// 	if err != nil {
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			if pgErr.Code == "23505" { // Unique constraint violation
// 				return 0, fmt.Errorf("%s: %w", op, storage.ErrUrlExists)
// 			}
// 		}

// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	// LastInsertId() - этот метод не работает с PostgreSQL
// 	return id, nil
// }

// func (s *Storage) GetUrl(alias string) (string, error) {
// 	const op = "storage.postgre.GetUrl"

// 	var url string
// 	query := "SELECT url FROM url WHERE alias = $1"

// 	err := s.db.QueryRow(query, alias).Scan(&url)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return "", fmt.Errorf("%s: %w", op, storage.ErrUrlNotFound)
// 		}
// 		return "", fmt.Errorf("%s: %w", op, err)
// 	}

// 	return url, nil
// }

// func (s *Storage) DeleteUrl(alias string) error {
// 	const op = "storage.postgre.DeleteUrl"

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	res, err := s.db.ExecContext(ctx, "DELETE FROM url WHERE alias = $1", alias)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("%s: %w", op, storage.ErrUrlNotFound)
// 	}

// 	return nil
// }
