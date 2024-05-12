package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	cahce_memory "github.com/maximpontryagin/level0/internal/storage/cahcememory"
)

func StartServer(c *cahce_memory.Cache) error {
	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		HandlerGetOnly(w, r, c)
	})
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return err
	}
	log.Println("Server started listening on port 8000")
	return nil
}

func HandlerGetOnly(w http.ResponseWriter, r *http.Request, c *cahce_memory.Cache) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	if r.Method == http.MethodGet {
		url := strings.Split(r.URL.Path, "/")
		if len(url) != 3 {
			http.Error(w, "Неверный формат строки", http.StatusBadRequest)
			log.Println(url)
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
