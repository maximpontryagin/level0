package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
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
	log.Println("✓ Таблицы успешно созданы")
	return nil
}
