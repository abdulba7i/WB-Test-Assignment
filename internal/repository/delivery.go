package repository

import (
	"database/sql"
	"fmt"
	"l0/internal/model"
)

func (s *Storage) AddDelivery(tx *sql.Tx, delivery model.Delivery) (int64, error) {
	const op = "storage.postgres.AddDelivery"

	var id int64
	query := "INSERT INTO delivery (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	err := tx.QueryRow(query, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
