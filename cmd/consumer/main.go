package main

import (
	"fmt"
	"log"

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

	fmt.Println("База данных подключилась")
	//Созданик таблицы в БД если они еще не созданы
	err = postgres.CreateDBtable(db)
	if err != nil {
		log.Println(err)
	}

	// Инициализация кеша
	cache := cahce_memory.New()
	// Заполнение кеша данными из БД (для случая отключения http сервера)
	cache.WritingCahce(db)

	go func() {
		err := nats_streaming.ConnectNats(db, cache)
		if err != nil {
			fmt.Println("Ошибка подключения к NATS:", err)
		}
	}()

	go func() {
		err = server.StartServer(cache)
		if err != nil {
			log.Println(err)
		}
	}()
	fmt.Println("http сервер запущен")

	select {}
}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config_example")
	return viper.ReadInConfig()
}
