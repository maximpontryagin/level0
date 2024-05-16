package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
)

func main() {
	// Подключение к NATS Streaming Server
	stanConn, err := stan.Connect("test-cluster", "publisher-client-id")
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS Streaming: %v", err)
	}
	defer stanConn.Close()

	// Сериализация json в данные структуры
	jsonData, err := os.ReadFile("model_data.json")
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	// Тема, в которую загрузится сообщение в NATS Streaming Server
	publisherSubject := "publisher_subject"
	err = stanConn.Publish(publisherSubject, jsonData)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	} else {
		log.Println("Данные отправлены в NATS Streaming:", string(jsonData))
	}
}
