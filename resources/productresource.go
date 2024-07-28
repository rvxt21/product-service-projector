package resources

import (
	"encoding/json"
	"errors"
	"net/http"
	"products/enteties"
	"products/middleware"
	"products/storage"

	"github.com/rs/zerolog/log"
)

type ProductsResourse struct {
	S *storage.Storage
}

func (tr *ProductsResourse) RegisterRoutes(m *http.ServeMux) {
	m.HandleFunc("POST /products", tr.CreateProduct)
	m.Handle("DELETE /products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.DeleteProduct)))
	m.Handle("PATCH /products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.UpdateAvailability)))
} //alternative for register routes

func (tr *ProductsResourse) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product enteties.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id := tr.S.CreateOneProduct(product)
	product.ID = id
	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (tr *ProductsResourse) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (tr *ProductsResourse) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.IdKey).(int)
	if tr.S.DeleteProduct(id) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Product not found", http.StatusNotFound)
	}
}

func (tr *ProductsResourse) UpdateProduct(w http.ResponseWriter, r *http.Request) {

}

type UpdateAvailabilityRequest struct {
	IsAvailable bool `json:"is_available"`
}

func (tr *ProductsResourse) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	var req UpdateAvailabilityRequest

	productID := r.Context().Value(middleware.IdKey).(int)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Err(err).Msg("error to update availability")
		return
	}

	err := tr.S.UpdateAvailability(productID, req.IsAvailable)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}
