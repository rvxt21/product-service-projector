package resources

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"products/enteties"
	"products/middleware"
	"products/storage"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type ProductsResourse struct {
	S *storage.DBStorage
}

func (tr *ProductsResourse) RegisterRoutes(r *mux.Router) {
	r.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.DeleteProduct))).Methods("DELETE")
	r.HandleFunc("/products/{id}", tr.UpdateProduct).Methods("PUT")
	r.Handle("/products/availability/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.UpdateAvailability))).Methods("PATCH")
	r.HandleFunc("/products", tr.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", tr.GetProductByID).Methods("GET")
	r.HandleFunc("/products", tr.GetAllProducts).Methods("GET")
	r.HandleFunc("/products/by-ids", tr.GetProductsByIDS).Methods("GET")
	//r.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.GetByID))).Methods("GET")

	r.HandleFunc("/categories", tr.CreateCategory).Methods("POST")
	r.HandleFunc("/categories", tr.GetAllCategories).Methods("GET")
	r.Handle("/categories/{idCategory}", middleware.IdMiddlewareCategory(http.HandlerFunc(tr.GetCategoryByID))).Methods("GET")
	r.Handle("/categories/{idCategory}", middleware.IdMiddlewareCategory(http.HandlerFunc(tr.UpdateCategory))).Methods("PUT")
	r.Handle("/categories/{idCategory}", middleware.IdMiddlewareCategory(http.HandlerFunc(tr.DeleteCategory))).Methods("DELETE")

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
	strLimit := r.URL.Query().Get("limit")
	strOffset := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(strLimit)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset, err := strconv.Atoi(strOffset)
	if err != nil || offset < 1 {
		offset = 0
	}
	products, err := tr.S.GetAllProductsDb(limit, offset)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (tr *ProductsResourse) GetProductsByIDS(w http.ResponseWriter, r *http.Request) {
	idsStr := r.URL.Query().Get("ids")
	if idsStr == "" {
		http.Error(w, "No IDs provided", http.StatusBadRequest)
		return
	}

	products, err := tr.S.GetProductsByIDSDB(idsStr)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Error().Msgf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (tr *ProductsResourse) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := tr.S.GetProductByIDDb(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

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

func (tr *ProductsResourse) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product enteties.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	product.ID = id

	if err := tr.S.UpdateProductBd(product); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

	err := tr.S.UpdateProductAvailabilityDB(productID, req.IsAvailable)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		http.Error(w, "Unable to update availability", http.StatusInternalServerError)
	}
}

func (tr *ProductsResourse) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category enteties.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := tr.S.CreateCategory(category)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create category in database")
		http.Error(w, "Unable to create category", http.StatusInternalServerError)
		return
	}
	category.IdCategory = id
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(category); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}

func (tr *ProductsResourse) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := tr.S.GetAllCategoriesDb()
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (tr *ProductsResourse) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.IdKey).(int)
	category, found, err := tr.S.GetCategoryByID(id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (tr *ProductsResourse) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var category enteties.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	category.IdCategory = r.Context().Value(middleware.IdKey).(int)
	err := tr.S.UpdateCategory(category)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Unable to update category", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (tr *ProductsResourse) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.IdKey).(int)
	success, err := tr.S.DeleteCategory(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete category from database")
		http.Error(w, "Unable to delete category", http.StatusInternalServerError)
		return
	}
	if success {
		w.WriteHeader(http.StatusNoContent)
	} else {
		log.Error().Msgf("Category with ID %d not found", id)
		http.Error(w, "Category not found", http.StatusNotFound)
	}
}
