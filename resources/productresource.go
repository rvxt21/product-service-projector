package resources

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"products/enteties"
	"products/middleware"
	"products/storage"

	"github.com/gorilla/mux"
)

type ProductsResourse struct {
	S *storage.DBStorage
}

func (tr *ProductsResourse) RegisterRoutes(r *mux.Router) {
	r.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.DeleteProduct))).Methods("DELETE")
	//r.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.UpdateProduct))).Methods("PUT")
	r.Handle("/products/availability/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.UpdateAvailability))).Methods("PATCH")
	r.HandleFunc("/products", tr.CreateProduct).Methods("POST")
	r.HandleFunc("/products", tr.GetAllProducts).Methods("GET")
	//r.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.GetByID))).Methods("GET")
} //alternative for register routes

func (tr *ProductsResourse) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product enteties.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := tr.S.CreateOneProductDb(product)
	if err != nil {
		log.Printf("Failed to create product in database: %v", err)
		http.Error(w, "Unable to create product", http.StatusInternalServerError)
		return
	}
	product.ID = id
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (tr *ProductsResourse) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := tr.S.GetAllProductsDb()

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// func (tr *ProductsResourse) GetByID(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	idStr := vars["id"]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid product ID", http.StatusBadRequest)
// 		return
// 	}

// 	product, found := tr.S.GetProductByID(id)
// 	if !found {
// 		http.Error(w, "Product not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(product); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		return
// 	}
// }

func (tr *ProductsResourse) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.IdKey).(int)
	success, err := tr.S.DeleteProductDb(id)
	if err != nil {
		log.Printf("Failed to delete product from database: %v", err)
		http.Error(w, "Unable to delete product", http.StatusInternalServerError)
		return
	}
	if success {
		w.WriteHeader(http.StatusNoContent)
	} else {
		log.Printf("Product with ID %d not found", id)
		http.Error(w, "Product not found", http.StatusNotFound)
	}
}

// func (tr *ProductsResourse) UpdateProduct(w http.ResponseWriter, r *http.Request) {
// 	var product enteties.Product
// 	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	product.ID = r.Context().Value(middleware.IdKey).(int)
// 	err := tr.S.UpdateProduct(product)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrProductNotFound) {
// 			http.Error(w, "Product not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, "Unable to update product", http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(product)
// }

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

	err := tr.S.UpdateProductAvailability(productID, req.IsAvailable)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		http.Error(w, "Unable to update availability", http.StatusInternalServerError)
	}
}
