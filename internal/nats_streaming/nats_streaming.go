package nats_streaming

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	cahce_memory "github.com/maximpontryagin/level0/internal/storage/cahcememory"
	"github.com/maximpontryagin/level0/internal/struct_delivery"
	"github.com/nats-io/nats.go"
)

func ConnectNats(db *sql.DB, c *cahce_memory.Cache) error {
	//подключение к nats серверу
	fmt.Println("подключаюсь к nats серверу")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	// подписка на канал

	publisher_subject := "publisher_subject"
	subscription, err := nc.Subscribe(publisher_subject, func(m *nats.Msg) {
		var order struct_delivery.Order
		log.Printf("успешно подключился к Nats streamig. Получено сообщение: %s\n", string(m.Data))

		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			log.Println(err)
		}

		// Записывание данных из nats server в БД
		query_order := `
			INSERT INTO orders (order_uid, track_number, entry, delivery_name, delivery_phone, delivery_zip, delivery_city, delivery_address, delivery_region, delivery_email, payment_transaction, payment_request_id, payment_currency, payment_provider, payment_amount, payment_dt, payment_bank, payment_delivery_cost, payment_goods_total, payment_custom_fee, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)
			RETURNING order_id
		`

		var orderID int

		err = db.QueryRow(query_order,
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
			}
		}
		// Записывание данных из nats server в кеш
		c.Set(strconv.Itoa(orderID), order)
		res_cache, _ := c.Get(strconv.Itoa(orderID))
		log.Println("В кеш записано сообщение:", res_cache)
	})
	if err != nil {
		log.Println(err)
	}
	defer subscription.Unsubscribe()

	// "Блокирование" программы, что бы она продолжала слушать nats server
	select {}

}
