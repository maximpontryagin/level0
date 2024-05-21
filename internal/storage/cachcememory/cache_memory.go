package cahce_memory

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/maximpontryagin/level0/internal/struct_delivery"
)

// Cache - универсальная структура кэша.
// T - тип хранимых значений.
type Cache[T any] struct {
	sync.RWMutex // Для безопасной записи данных при вызове go рутин
	store        map[string]T
}

// Инициализация хранилища
func New[T any]() *Cache[T] {
	return &Cache[T]{store: make(map[string]T)}
}

// Установливает значения в кеш
func (c *Cache[T]) Set(key string, value T) {
	c.Lock()
	defer c.Unlock()
	c.store[key] = value
}

// Возвращает значение из кеша
func (c *Cache[T]) Get(key string) (T, bool) {
	c.RLock()
	defer c.RUnlock()

	val, found := c.store[key]
	return val, found
}

// Удаление кеша по переданному ключу
func (c *Cache[T]) Delete(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, found := c.store[key]; !found {
		return errors.New("ключ в кеше не найден")
	}
	delete(c.store, key)
	return nil
}

// Записывание данных из БД в кеш
func Writing_In_Cahce_from_DB(c *Cache[struct_delivery.Order], db *sql.DB) error {
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
			log.Printf("error scanning orders row: %v", err)
			continue
		}
		c.Set(strconv.Itoa(orderID), order)
	}
	log.Println("Данные из БД записаны в кэш")
	if err := ordersRows.Err(); err != nil {
		return err
	}
	return nil
}
