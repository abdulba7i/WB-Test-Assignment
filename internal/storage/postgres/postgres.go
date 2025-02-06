package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"l0wb/internal/config"
	"l0wb/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SMID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OOFShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func New(c config.Database) (*Storage, error) {
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
		order_uid VARCHAR(255) NOT NULL,
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

func (s *Storage) AddDelivery(delivery Delivery) (int64, error) {
	const op = "storage.postgres.AddDelivery"

	var id int64
	query := "INSERT INTO delivery (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	err := s.db.QueryRow(query, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) AddPayment(payment Payment) (int64, error) {
	const op = "storage.postgres.AddPayment"

	var id int64
	query := "INSERT INTO payment (transaction, request_id, currency, provider, amount, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	err := s.db.QueryRow(query, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) AddItems(order_uid string, items []Item) error {
	const op = "storage.postgres.AddItems"

	query := "INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	for _, item := range items {
		err := s.db.QueryRow(query, order_uid, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) AddOrder(ordr Order) error {
	var err error
	const op = "storage.postgres.AddOrder"
	tx, err := s.db.Begin()
	defer tx.Rollback()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	idDvr, errD := s.AddDelivery(ordr.Delivery)
	if errD != nil {
		return fmt.Errorf("%s: %w", op, errD)
	}

	idPymnt, errP := s.AddPayment(ordr.Payment)
	if errP != nil {
		return fmt.Errorf("%s: %w", op, errP)
	}

	query := `INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, oof_shard)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err = s.db.Exec(query, ordr.OrderUID, ordr.TrackNumber, ordr.Entry, idDvr, idPymnt, ordr.Locale, ordr.InternalSignature, ordr.CustomerID, ordr.DeliveryService, ordr.ShardKey, ordr.SMID, ordr.OOFShard)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Теперь можно добавлять items, так как order_uid уже в orders
	errI := s.AddItems(ordr.OrderUID, ordr.Items)
	if errI != nil {
		return fmt.Errorf("%s: %w", op, errI)
	}
	tx.Commit()

	return nil
}

func (s *Storage) GetOrderById(id string) (Order, error) {
	const op = "storage.postgres.GetOrderById"

	query := `SELECT * FROM orders WHERE order_uid = $1 JOIN delivery on orders.delivery_id = delivery.id join payment on orders.payment_id = payment.id JOIN items on orders.order_uid = items.order_uid`
	var ordr Order

	err := s.db.QueryRow(query, id).Scan(&ordr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, fmt.Errorf("%s: %w", op, storage.ErrUrlNotFound)
		}
		return Order{}, fmt.Errorf("%s: %w", op, err)
	}

	return ordr, nil
}
