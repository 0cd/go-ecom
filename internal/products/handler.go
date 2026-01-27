package products

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

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Printf("Failed to list products: %v", err)
		http.Error(w, "failed to retrieve products", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) FindProductByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Invalid product id: %s: %v", idParam, err)
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.service.FindProductByID(r.Context(), id)
	if err != nil {
		log.Printf("Failed to find product %d: %v", id, err)
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	json.Write(w, http.StatusOK, product)
}
