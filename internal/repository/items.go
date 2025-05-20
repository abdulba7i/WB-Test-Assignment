package repository

import (
	"database/sql"
	"fmt"
	"l0/internal/model"
)

func (s *Storage) AddItems(tx *sql.Tx, order_uid string, items []model.Item) error {
	const op = "storage.postgres.AddItems"

	query := "INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	for _, item := range items {
		_, err := tx.Exec(query, order_uid, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
