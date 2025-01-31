package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
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
		delivery_id INTEGER NOT NULL,
		payment_id INTEGER NOT NULL,
		locale VARCHAR(255) NOT NULL,
		internal_signature VARCHAR(255) NOT NULL,
		customer_id VARCHAR(255) NOT NULL,
		delivery_service VARCHAR(255) NOT NULL,
		shardkey VARCHAR(255) NOT NULL,
		sm_id INTEGER NOT NULL,
		oof_shard VARCHAR(255) NOT NULL,
		date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS delivery (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		phone VARCHAR(20) NOT NULL,
		zip VARCHAR(20) NOT NULL,
		city VARCHAR(255) NOT NULL,
		address VARCHAR(255) NOT NULL,
		region VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS payment (
		id SERIAL PRIMARY KEY,
		transaction VARCHAR(255) NOT NULL,
		request_id VARCHAR(255) NOT NULL,
		currency VARCHAR(10) NOT NULL,
		provider VARCHAR(50) NOT NULL,
		amount INTEGER NOT NULL,
		bank VARCHAR(255) NOT NULL,
		delivery_cost INTEGER NOT NULL,
		goods_total INTEGER NOT NULL,
		custom_fee INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		chrt_id INTEGER NOT NULL,
		track_number VARCHAR(255) NOT NULL,
		price INTEGER NOT NULL,
		rid VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		sale INTEGER NOT NULL,
		size VARCHAR(255) NOT NULL,
		total_price INTEGER NOT NULL,
		nm_id INTEGER NOT NULL,
		brand VARCHAR(255) NOT NULL,
		status INTEGER NOT NULL
	);
	`

	_, err = db.Exec(stmt1)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// created index
	stmt := `
	CREATE INDEX IF NOT EXISTS order_uid_idx ON orders (order_uid);
	CREATE INDEX IF NOT EXISTS delivery_id_idx ON orders (delivery_id);
	CREATE INDEX IF NOT EXISTS payment_id_idx ON orders (payment_id);
	CREATE INDEX IF NOT EXISTS items_idx ON items (rid);
	`

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{db: db}, nil
}
