package main

import (
	"encoding/json"
	"net/http"

	"github.com/andrereliquias/gointensivo/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	//chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/order", Order)
	// http.HandleFunc("/order", Order)
	http.ListenAndServe(":8888", r)
}

func Order(w http.ResponseWriter, r *http.Request) {
	order, err := entity.NewOrder("1234", 1000, 10)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	order.CalculateFinalPrice()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(order)

}
