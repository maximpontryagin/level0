package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/maximpontryagin/level0/internal/struct_delivery"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func ConnectDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateDBtable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS orders (
		order_id SERIAL PRIMARY KEY,
		order_uid VARCHAR(255) UNIQUE NOT NULL,
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
		payment_currency VARCHAR(5) NOT NULL,
		payment_provider VARCHAR(255) NOT NULL,
		payment_amount INT NOT NULL,
		payment_dt INT NOT NULL,
		payment_bank VARCHAR(255) NOT NULL,
		payment_delivery_cost INT NOT NULL,
		payment_goods_total INT NOT NULL,
		payment_custom_fee INT NOT NULL,
		locale VARCHAR(10),
		internal_signature VARCHAR(255),
		customer_id VARCHAR(255) NOT NULL,
		delivery_service VARCHAR(255),
		shardkey VARCHAR(255) NOT NULL,
		sm_id INT,
		date_created TIMESTAMP NOT NULL,
		oof_shard VARCHAR(255) NOT NULL
	);
		CREATE TABLE IF NOT EXISTS order_items (
		item_id SERIAL PRIMARY KEY,
		order_id INT,
		chrt_id INT NOT NULL,
		track_number VARCHAR(255) NOT NULL,
		price INT NOT NULL,
		rid VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		sale INT NOT NULL,
		size VARCHAR(255) NOT NULL,
		total_price INT NOT NULL,
		nm_id INT NOT NULL,
		brand VARCHAR(255) NOT NULL,
		status INT NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func Writing_in_DB(db *sql.DB, order struct_delivery.Order) (int, error) {
	query_order := `
            INSERT INTO orders (order_uid, track_number, entry, delivery_name, delivery_phone, delivery_zip, delivery_city, delivery_address, delivery_region, delivery_email, payment_transaction, payment_request_id, payment_currency, payment_provider, payment_amount, payment_dt, payment_bank, payment_delivery_cost, payment_goods_total, payment_custom_fee, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)
            RETURNING order_id
        `

	var orderID int
	err := db.QueryRow(query_order,
		order.OrderUID, order.TrackNumber, order.Entry,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region,
		order.Delivery.Email, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee, order.Locale,
		order.InternalSign, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SMID, order.DateCreated, order.OOFShard).Scan(&orderID)

	if err != nil {
		log.Println(err)
		return orderID, err
	}

	for _, item := range order.Items {
		query_item := `
            INSERT INTO order_items (order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        `
		_, err := db.Exec(query_item, orderID, item.ChrtID, item.TrackNumber, item.Price, item.RID,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			log.Println(err)
			return orderID, err
		}
	}
	return orderID, nil

}
