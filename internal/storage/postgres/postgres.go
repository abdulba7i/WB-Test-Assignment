package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

type Storage struct {
	db *sql.DB
}
type DB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

func New(c DB) (*Storage, error) {
	const op = "storage.postgre.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	// created table
	stmt1 := `
	CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(255) NOT NULL,
    delivery_name VARCHAR(255) NOT NULL,
    delivery_phone VARCHAR(20) NOT NULL,
    delivery_zip VARCHAR(20) NOT NULL,
    delivery_city VARCHAR(255) NOT NULL,
    delivery_address VARCHAR(255) NOT NULL,
    delivery_region VARCHAR(255) NOT NULL,
    delivery_email VARCHAR(255) NOT NULL,
    payment_transaction VARCHAR(255) NOT NULL,
    payment_request_id VARCHAR(255) NOT NULL,
    payment_currency VARCHAR(10) NOT NULL,
    payment_provider VARCHAR(50) NOT NULL,
    payment_amount INTEGER NOT NULL,
    payment_bank VARCHAR(255) NOT NULL,
    payment_delivery_cost INTEGER NOT NULL,
    payment_goods_total INTEGER NOT NULL,
    payment_custom_fee INTEGER NOT NULL,
    date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP);`

	_, err = db.Exec(stmt1)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// created index
	stmt := `
	CREATE INDEX IF NOT EXISTS idx_track_number ON orders(track_number);
	CREATE INDEX IF NOT EXISTS idx_date_created ON orders(date_created);
	CREATE INDEX IF NOT EXISTS idx_customer_id ON orders(customer_id);
	CREATE INDEX IF NOT EXISTS idx_payment_transaction ON orders(payment_transaction);
	CREATE INDEX IF NOT EXISTS idx_delivery_phone ON orders(delivery_phone);
	`

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{db: db}, nil
}
