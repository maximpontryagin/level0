package nats_streaming

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	cachce_memory "github.com/maximpontryagin/level0/internal/storage/cachcememory"
	"github.com/maximpontryagin/level0/internal/storage/postgres"
	"github.com/maximpontryagin/level0/internal/struct_delivery"
	"github.com/nats-io/stan.go"
)

func DurableSubscriptions(stanConn stan.Conn, db *sql.DB, c *cachce_memory.Cache) (stan.Subscription, error) {
	subscribe, err := stanConn.Subscribe("publisher_subject", func(m *stan.Msg) {
		var order struct_delivery.Order
		log.Printf("Получено сообщение из nats streaming server: %s\n", string(m.Data))

		// Сериализация данных в структуру order
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			log.Println(err)
			return
		}

		// Записывание данных из NATS server в БД
		orderID, err := postgres.Writing_in_DB(db, order)
		if err != nil {
			log.Println(err)
			return
		}

		// Записывание данных из NATS server в кеш
		c.Set(strconv.Itoa(orderID), order)

	}, stan.DurableName("my-durable-subscription"))
	if err != nil {
		return nil, err
	}
	return subscribe, nil
}
