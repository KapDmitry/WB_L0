package repo

import (
	"context"
	"database/sql"

	"github.com/KapDmitry/WB_L0/internal/order"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(d *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		DB: d,
	}
}

func (p *PostgresRepo) Add(ctx context.Context, ord order.Order) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = insertOrdersList(tx, ord)
	if err != nil {
		return err
	}

	err = insertDelivery(tx, ord.OrderUID, ord.OrderDelivery)
	if err != nil {
		return err
	}

	err = insertPayment(tx, ord.OrderUID, ord.OrderPayment)
	if err != nil {
		return err
	}

	err = insertItems(tx, ord.OrderUID, ord.Items)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func insertOrdersList(tx *sql.Tx, orderData order.Order) error {
	query := `
		INSERT INTO orderslist (order_uid, track_number, order_entry, locale, internal_signature, customer_id, 
			delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	_, err := tx.Exec(query, orderData.OrderUID, orderData.TrackNumber, orderData.Entry,
		orderData.Locale, orderData.InternalSignature, orderData.CustomerID, orderData.DeliveryService,
		orderData.Shardkey, orderData.SMID, orderData.DateCreated, orderData.OOFShard)

	return err
}

func insertDelivery(tx *sql.Tx, orderID string, deliveryData order.Delivery) error {
	query := `
		INSERT INTO delivery (order_id, delivery_name, phone, zip, city, delivery_address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.Exec(query, orderID, deliveryData.Name, deliveryData.Phone, deliveryData.ZIP,
		deliveryData.City, deliveryData.Address, deliveryData.Region, deliveryData.Email)

	return err
}

func insertPayment(tx *sql.Tx, orderID string, paymentData order.Payment) error {
	query := `
		INSERT INTO payment (order_id, payment_transaction, request_id, currency, payment_provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := tx.Exec(query, orderID, paymentData.Transaction, paymentData.RequestID,
		paymentData.Currency, paymentData.Provider, paymentData.Amount, paymentData.PaymentDT,
		paymentData.Bank, paymentData.DeliveryCost, paymentData.GoodsTotal, paymentData.CustomFee)

	return err
}

func insertItems(tx *sql.Tx, orderID string, itemsData []order.Item) error {
	query := `
		INSERT INTO items (order_id, chrt_id, track_number, price, rid, name_item, sale,
			size_item, total_price, nm_id, brand, status_item)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range itemsData {
		_, err := tx.Exec(query, orderID, item.ChrtID, item.TrackNumber, item.Price,
			item.RID, item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NMID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresRepo) getItems(ctx context.Context, id string) ([]order.Item, error) {
	rows, err := p.DB.QueryContext(ctx, "SELECT chrt_id, track_number, price, rid, name_item, "+
		"sale, size_item, total_price, nm_id, brand, status_item FROM items WHERE order_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []order.Item
	for rows.Next() {
		var item order.Item
		err = rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NMID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (p *PostgresRepo) GetAll(ctx context.Context) ([]order.Order, error) {
	rows, err := p.DB.QueryContext(ctx, "SELECT orderslist.order_uid, orderslist.track_number, orderslist.order_entry, "+
		"orderslist.locale, orderslist.internal_signature, orderslist.customer_id, "+
		"orderslist.delivery_service, orderslist.shardkey, orderslist.sm_id, "+
		"orderslist.date_created, orderslist.oof_shard, "+
		"delivery.delivery_name, delivery.phone, delivery.zip, "+
		"delivery.city, delivery.delivery_address, delivery.region, delivery.email, "+
		"payment.payment_transaction, payment.request_id, payment.currency, "+
		"payment.payment_provider, payment.amount, payment.payment_dt, "+
		"payment.bank, payment.delivery_cost, payment.goods_total, payment.custom_fee "+
		"FROM orderslist "+
		"JOIN delivery ON orderslist.order_uid = delivery.order_id "+
		"JOIN payment ON orderslist.order_uid = payment.order_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]order.Order, 0)
	for rows.Next() {
		var curOrder order.Order
		err := rows.Scan(
			&curOrder.OrderUID, &curOrder.TrackNumber, &curOrder.Entry,
			&curOrder.Locale, &curOrder.InternalSignature, &curOrder.CustomerID,
			&curOrder.DeliveryService, &curOrder.Shardkey, &curOrder.SMID,
			&curOrder.DateCreated, &curOrder.OOFShard,
			&curOrder.OrderDelivery.Name, &curOrder.OrderDelivery.Phone, &curOrder.OrderDelivery.ZIP,
			&curOrder.OrderDelivery.City, &curOrder.OrderDelivery.Address, &curOrder.OrderDelivery.Region, &curOrder.OrderDelivery.Email,
			&curOrder.OrderPayment.Transaction, &curOrder.OrderPayment.RequestID, &curOrder.OrderPayment.Currency,
			&curOrder.OrderPayment.Provider, &curOrder.OrderPayment.Amount, &curOrder.OrderPayment.PaymentDT,
			&curOrder.OrderPayment.Bank, &curOrder.OrderPayment.DeliveryCost, &curOrder.OrderPayment.GoodsTotal, &curOrder.OrderPayment.CustomFee,
		)
		if err != nil {
			return nil, err
		}

		curOrder.Items, err = p.getItems(ctx, curOrder.OrderUID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, curOrder)

	}

	return orders, nil

}
