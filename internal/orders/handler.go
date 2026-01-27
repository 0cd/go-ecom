package orders

import (
	"log"
	"net/http"

	"github.com/0cd/go-ecom/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var order createOrderParams
	if err := json.Read(r, &order); err != nil {
		log.Printf("Failed to parse order request: %v", err)
		http.Error(w, "invalid order request body", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), order)
	if err != nil {
		log.Printf("Failed to create order in service: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)
}
