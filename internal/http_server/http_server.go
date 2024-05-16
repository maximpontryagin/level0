package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	cachce_memory "github.com/maximpontryagin/level0/internal/storage/cachcememory"
)

func StartServer(c *cachce_memory.Cache, ctx context.Context) error {
	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		HandlerGetOnly(w, r, c)
	})

	httpServer := &http.Server{
		Addr: ":8000",
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()
	log.Println("Сервер запустился и слушает порт", httpServer.Addr)
	// Ждём сигнал о завершении работы сервиса
	<-ctx.Done()
	// Завершаем работу http сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second)
	defer shutdownCancel()
	return httpServer.Shutdown(shutdownCtx)
}

func HandlerGetOnly(w http.ResponseWriter, r *http.Request, c *cachce_memory.Cache) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	if r.Method == http.MethodGet {
		url := strings.Split(r.URL.Path, "/")
		if len(url) != 3 {
			http.Error(w, "Неверный формат строки", http.StatusBadRequest)
			return
		}
		order_id := url[2]
		_, err := strconv.Atoi(order_id)
		if err != nil {
			http.Error(w, "Неверный формат id заказа", http.StatusNotFound)
			return
		}
		order, search_result := c.Get(order_id)
		if !search_result {
			http.Error(w, "Введен несуществующий id заказа", http.StatusNotFound)
			return
		}
		err = json.NewEncoder(w).Encode(order)
		if err != nil {
			http.Error(w, "Ошибка в отправке Json", http.StatusInternalServerError)
			return
		}
	}
}
