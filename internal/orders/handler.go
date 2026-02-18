package orders

import (
	"log"
	"net/http"
	"strconv"

	"github.com/0cd/go-ecom/internal/json"
	"github.com/go-chi/chi/v5"
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
	userID, ok := r.Context().Value("userID").(int64)
	if !ok || userID == 0 {
		log.Printf("Failed to extract userID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	order := createOrderParams{
		UserID: userID,
	}
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

func (h *handler) FindOrderByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Invalid order id: %s: %v", idParam, err)
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.service.FindOrderByID(r.Context(), id)
	if err != nil {
		log.Printf("Failed to find order %d: %v", id, err)
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	json.Write(w, http.StatusOK, order)
}
