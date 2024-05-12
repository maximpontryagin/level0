package cahce_memory

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/maximpontryagin/level0/internal/struct_delivery"
)

// Струкрута cache описывает хранилище
type Cache struct {
	sync.RWMutex // Для безопасной записи данных при вызове go рутин
	items        map[string]Item
}

type Item struct {
	Value interface{} // Любое значение
}

// Инициализация хранилища
func New() *Cache {
	items := make(map[string]Item) // инициализируем карту в паре ключ(string)/значение(Item)
	cache := Cache{
		items: items,
	}
	return &cache
}

// Установка значений
func (c *Cache) Set(key string, value interface{}) {

	c.Lock()
	defer c.Unlock()

	c.items[key] = Item{
		Value: value,
	}
}

// Получение значений
func (c *Cache) Get(key string) (interface{}, bool) {

	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	return item.Value, true
}

// Удаление кеша по переданному ключу
func (c *Cache) Delete(key string) error {

	c.Lock()
	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("ключ в кеше не найден")
	}
	delete(c.items, key)

	return nil
}

// Записывание данных из БД в кеш
func (c *Cache) WritingCahce(db *sql.DB) error {
	ordersQuery := "SELECT * FROM orders"
	ordersRows, err := db.Query(ordersQuery)
	if err != nil {
		return err
	}
	defer ordersRows.Close()

	for ordersRows.Next() {
		var orderID int
		var order struct_delivery.Order

		err := ordersRows.Scan(&orderID,
			&order.OrderUID, &order.TrackNumber, &order.Entry,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
			&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
			&order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestID,
			&order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
			&order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee, &order.Locale,
			&order.InternalSign, &order.CustomerID, &order.DeliveryService,
			&order.ShardKey, &order.SMID, &order.DateCreated, &order.OOFShard)
		if err != nil {
			log.Printf("erorr scanning orders row: %v", err)
			continue
		}
		c.Set(strconv.Itoa(orderID), order)
		cache_res, _ := c.Get(strconv.Itoa(orderID))
		log.Println("Данные записанные в кеш:", cache_res)
	}
	if err := ordersRows.Err(); err != nil {
		return err
	}
	return nil
}
