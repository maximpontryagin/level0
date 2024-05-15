package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	server "github.com/maximpontryagin/level0/internal/http_server"
	"github.com/maximpontryagin/level0/internal/nats_streaming"
	cachce_memory "github.com/maximpontryagin/level0/internal/storage/cachcememory"
	"github.com/maximpontryagin/level0/internal/storage/postgres"
	"github.com/nats-io/nats.go"
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
	cache := cachce_memory.New()
	// Заполнение кеша данными из БД (для случая отключения http сервера)
	err = cache.Writing_In_Cahce_from_DB(db)
	if err != nil {
		log.Println("Ошибка заполнения кеша данными из бд:", err)
	}
	// Запуск nats streaming
	go func() {
		log.Println("Подключаюсь к nats серверу...")
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			log.Println("Ошибка подключения к NATS:", err)
			return
		}
		log.Println("Nats сервер подключен")
		subscription := nats_streaming.Writing_in_DB_and_Cache(db, cache, nc)
		defer subscription.Unsubscribe()
		select {}
	}()

	// Запуск сервера
	go func() {
		err = server.StartServer(cache)
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("http сервер запущен")
	// "Блокирование" go рутины main
	FinishSignal := make(chan os.Signal, 1)
	signal.Notify(FinishSignal, os.Interrupt)
	signal.Notify(FinishSignal, syscall.SIGTERM)

	<-FinishSignal
	log.Println("Сервис завершает работу и закрывает соединия с Nats и БД...")
}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config_example")
	return viper.ReadInConfig()
}
