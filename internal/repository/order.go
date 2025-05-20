package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"l0/internal/model"
	"log"
)

func (s *Storage) AddOrder(ordr model.Order) error {
	var err error
	const op = "storage.postgres.AddOrder"
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			log.Println("ROllback")
			tx.Rollback()
		} else {
			log.Println("Commit")
			tx.Commit()
		}
	}()

	idDvr, err := s.AddDelivery(tx, ordr.Delivery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	idPymnt, err := s.AddPayment(tx, ordr.Payment)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := `INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, oof_shard)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err = tx.Exec(query, ordr.OrderUID, ordr.TrackNumber, ordr.Entry, idDvr, idPymnt, ordr.Locale, ordr.InternalSignature, ordr.CustomerID, ordr.DeliveryService, ordr.ShardKey, ordr.SMID, ordr.OOFShard)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Теперь можно добавлять items, так как order_uid уже в orders
	err = s.AddItems(tx, ordr.OrderUID, ordr.Items)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetOrderById(id string) (model.Order, error) {
	const op = "storage.postgres.GetOrderById"

	var order model.Order
	var delivery model.Delivery
	var payment model.Payment
	var itemsJSON json.RawMessage

	query := `
		SELECT o.order_uid, o.track_number, o.entry, 
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email, 
			p.transaction, p.request_id, p.currency, p.provider, p.amount, 
			p.bank, p.delivery_cost, p.goods_total, p.custom_fee, 
			o.locale, o.internal_signature, o.customer_id, o.delivery_service, 
			o.shardkey, o.sm_id, o.date_created, o.oof_shard,
			COALESCE(i.items, '[]'::json) AS items
		FROM orders o
		JOIN delivery d ON o.delivery_id = d.id
		JOIN payment p ON o.payment_id = p.id
		LEFT JOIN LATERAL (
			SELECT json_agg(i) AS items
			FROM items i
			WHERE i.order_uid = o.order_uid
		) i ON true
		WHERE o.order_uid = $1
	`

	err := s.db.QueryRow(query, id).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
		&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount,
		&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
		&order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SMID, &order.DateCreated, &order.OOFShard,
		&itemsJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, fmt.Errorf("%s: order with id %s not found", op, id)
		}
		return model.Order{}, fmt.Errorf("%s: %w", op, err)
	}

	order.Delivery = delivery
	order.Payment = payment

	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return model.Order{}, fmt.Errorf("%s: failed to parse items JSON: %w", op, err)
	}

	return order, nil
}

func (s *Storage) GetAllOrders(limit, offset int) ([]model.Order, error) {
	const op = "storage.postgres.GetAllOrders"

	var orders []model.Order

	query := `
		SELECT o.order_uid, o.track_number, o.entry, 
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email, 
			p.transaction, p.request_id, p.currency, p.provider, p.amount, 
			p.bank, p.delivery_cost, p.goods_total, p.custom_fee, 
			o.locale, o.internal_signature, o.customer_id, o.delivery_service, 
			o.shardkey, o.sm_id, o.date_created, o.oof_shard,
			COALESCE(i.items, '[]'::json) AS items
		FROM orders o
		JOIN delivery d ON o.delivery_id = d.id
		JOIN payment p ON o.payment_id = p.id
		LEFT JOIN LATERAL (
			SELECT json_agg(i) AS items
			FROM items i
			WHERE i.order_uid = o.order_uid
		) i ON true
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, limit, offset)
	fmt.Print("LIMIT", limit, "OFFSET", offset)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var order model.Order
		var delivery model.Delivery
		var payment model.Payment
		var itemsJSON json.RawMessage

		if err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
			&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount,
			&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
			&order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.ShardKey, &order.SMID, &order.DateCreated, &order.OOFShard,
			&itemsJSON,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		order.Delivery = delivery
		order.Payment = payment

		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, fmt.Errorf("%s: failed to parse items JSON: %w", op, err)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// GetOrdersBatch получает заказы пакетами
func (s *Storage) GetOrdersBatch(ctx context.Context, batchSize int, processBatch func([]model.Order) error) error {
	const op = "storage.postgres.GetOrdersBatch"

	// Получаем общее количество заказов
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&total)
	if err != nil {
		return fmt.Errorf("%s: failed to get total count: %w", op, err)
	}

	// Обрабатываем данные пакетами
	for offset := 0; offset < total; offset += batchSize {
		query := `
			SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
				   o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
				   d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
				   p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt,
				   p.bank, p.delivery_cost, p.goods_total, p.custom_fee
			FROM orders o
			JOIN delivery d ON o.delivery_id = d.id
			JOIN payment p ON o.payment_id = p.id
			ORDER BY o.order_uid
			LIMIT $1 OFFSET $2`

		rows, err := s.db.QueryContext(ctx, query, batchSize, offset)
		if err != nil {
			return fmt.Errorf("%s: failed to query orders: %w", op, err)
		}
		defer rows.Close()

		var orders []model.Order
		for rows.Next() {
			var o model.Order
			var d model.Delivery
			var p model.Payment

			err := rows.Scan(
				&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
				&o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SMID, &o.DateCreated, &o.OOFShard,
				&d.Name, &d.Phone, &d.Zip, &d.City, &d.Address, &d.Region, &d.Email,
				&p.Transaction, &p.RequestID, &p.Currency, &p.Provider, &p.Amount, &p.PaymentDT,
				&p.Bank, &p.DeliveryCost, &p.GoodsTotal, &p.CustomFee,
			)
			if err != nil {
				return fmt.Errorf("%s: failed to scan order: %w", op, err)
			}

			// Получаем items для заказа
			items, err := s.getItemsForOrder(ctx, o.OrderUID)
			if err != nil {
				return fmt.Errorf("%s: failed to get items for order %s: %w", op, o.OrderUID, err)
			}

			o.Delivery = d
			o.Payment = p
			o.Items = items
			orders = append(orders, o)
		}

		if err = rows.Err(); err != nil {
			return fmt.Errorf("%s: rows error: %w", op, err)
		}

		// Обрабатываем пакет заказов
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := processBatch(orders); err != nil {
				return fmt.Errorf("%s: failed to process batch: %w", op, err)
			}
		}
	}

	return nil
}

// getItemsForOrder получает все items для заказа
func (s *Storage) getItemsForOrder(ctx context.Context, orderUID string) ([]model.Item, error) {
	query := `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1`

	rows, err := s.db.QueryContext(ctx, query, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NMID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
