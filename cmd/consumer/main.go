package main

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
	server "github.com/maximpontryagin/level0/internal/http_server"
	"github.com/maximpontryagin/level0/internal/nats_streaming"
	cahce_memory "github.com/maximpontryagin/level0/internal/storage/cahcememory"
	"github.com/maximpontryagin/level0/internal/storage/postgres"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Println(err)
	}
	// Подключение к БД
	db, err := postgres.ConnectDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	log.Println("База данных подключилась")

	// Создаст таблицы в БД если они еще не созданы
	err = postgres.CreateDBtable(db)
	if err != nil {
		log.Println(err)
	}

	// Инициализация кеша
	cache := cahce_memory.New()
	// Заполнение кеша данными из БД (для случая отключения http сервера)
	cache.WritingCahce(db)

	// Mutex для гарантии записи данных и подключаем nats
	var mu sync.Mutex
	go func() {
		mu.Lock()
		err := nats_streaming.ConnectNats(db, cache)
		mu.Unlock()
		if err != nil {
			log.Println("Ошибка подключения к NATS:", err)
		}
	}()

	// Запуск сервера
	go func() {
		err = server.StartServer(cache)
		if err != nil {
			log.Println(err)
		}
	}()
	fmt.Println("http сервер запущен")
	// "Блокирование" go рутины main
	select {}
}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config_example")
	return viper.ReadInConfig()
}
