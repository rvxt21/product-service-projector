package resources

import (
	"encoding/json"
	"errors"
	"net/http"
	"products/enteties"
	"products/middleware"
	"products/storage"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type ProductsResourse struct {
	S *storage.Storage
}

func (tr *ProductsResourse) RegisterRoutes(m *http.ServeMux) {
	m.HandleFunc("POST /products", (tr.CreateProduct))
	m.Handle("DELETE /products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.DeleteProduct)))
	m.Handle("PATCH /products/availability/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.UpdateAvailability)))
	m.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.UpdateProduct)))
	m.HandleFunc("GET /products", tr.GetAll)
	m.Handle("GET /products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.GetByID)))
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
	var catalogue enteties.Catalogue
	response := struct {
		Products map[string]string
	}{
		Products: make(map[string]string),
	}
	for id, product := range catalogue.Products {
		response.Products[id] = product.Name
	}
	json.NewEncoder(w).Encode(response)
}

func (tr *ProductsResourse) GetByID(w http.ResponseWriter, r *http.Request) {
	var catalogue enteties.Catalogue
	vars := mux.Vars(r)
	id := vars["id"]

	product, exists := catalogue.Products[id]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	response := struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
		Category    string  `json:"category"`
		IsAvailable bool    `json:"is_available"`
	}{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Category:    product.Category,
		IsAvailable: product.IsAvailable,
	}
	json.NewEncoder(w).Encode(response)
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
	var product enteties.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	product.ID = r.Context().Value(middleware.IdKey).(int)
	err := tr.S.UpdateProduct(product)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Unable to update product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

type UpdateAvailabilityRequest struct {
	IsAvailable bool `json:"is_available"`
}

func (tr *ProductsResourse) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	var req UpdateAvailabilityRequest

	productID := r.Context().Value(middleware.IdKey).(int)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Err(err).Msg("Error to update availability")
		return
	}

	err := tr.S.UpdateAvailability(productID, req.IsAvailable)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		http.Error(w, "Unable to update availability", http.StatusInternalServerError)
	}
}
