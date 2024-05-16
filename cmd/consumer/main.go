package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
	server "github.com/maximpontryagin/level0/internal/http_server"
	"github.com/maximpontryagin/level0/internal/nats_streaming"
	cachce_memory "github.com/maximpontryagin/level0/internal/storage/cachcememory"
	"github.com/maximpontryagin/level0/internal/storage/postgres"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
)

func main() {
	// Правильное завершение работы сервиса
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		log.Println("Сервис завершает работу и закрывает соединия с Nats, БД и отключает http сервис...")
		cancel()
	}()

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
	defer func() {
		db.Close()
		log.Println("База данных закрылась")
	}()
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

	var wg sync.WaitGroup

	// Запуск nats streaming server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Подключаюсь к NATS Streaming Server...")
		stanConn, err := stan.Connect("test-cluster", "subscriber-client-id")
		if err != nil {
			log.Println(err)
		}
		subscribe, err := nats_streaming.DurableSubscriptions(stanConn, db, cache)
		if err != nil {
			log.Fatalf("Ошибка при подписке NATS Streaming Server: %v", err)
		}
		<-ctx.Done()

		// Отписываемся и закрываем nats streaming server
		err = subscribe.Unsubscribe()
		if err != nil {
			log.Printf("Ошибка отписки от nats streaming server : %s\n", err.Error())
		} else {
			log.Println("Успешно отписался от nats streaming server")
		}
		err = stanConn.Close()
		if err != nil {
			log.Printf("Ошибка закрытия nats streaming server : %s\n", err.Error())
		} else {
			log.Println("Успешно закрыл nats streaming server")
		}
	}()

	// Запуск сервера (Graceful shutdown внутри функции StartServer)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.StartServer(cache, ctx); err != nil {
			log.Println("Ошибка при работе HTTP сервера:", err)
		} else {
			log.Println("Успешно закрыл HTTP сервер")
		}
	}()
	wg.Wait()
}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config_example")
	return viper.ReadInConfig()
}
