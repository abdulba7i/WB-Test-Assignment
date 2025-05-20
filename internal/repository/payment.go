package repository

import (
	"database/sql"
	"fmt"
	"l0/internal/model"
)

func (s *Storage) AddPayment(tx *sql.Tx, payment model.Payment) (int64, error) {
	const op = "storage.postgres.AddPayment"

	var id int64
	query := "INSERT INTO payment (transaction, request_id, currency, provider, amount, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	err := tx.QueryRow(query, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
