package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

// publisher отправляет данные из json к брокеру (необходим для проверки работы)
func main() {
	// подключение к nats streaming
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Printf("Ошибка при подключению к nats: %v", err)
	}

	// Сериализация json в данные струкур
	json_data, err := os.ReadFile("model_data.json")
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
	}

	publisher_subject := "publisher_subject" // тема в которую загрузится сообщение в nats server
	nc.Publish(publisher_subject, json_data)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	} else {
		log.Println("Данные отправлены в NATS Streaming", string(json_data))
	}
}
