package resources

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"products/internal/enteties"
	"products/internal/middleware"
	"products/internal/storage"
	"products/pkg/utils"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (tr *ProductsResourse) GetProducts(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		tr.GetByName(w, r)
		return
	}

	tr.GetAllProducts(w, r)
}

func (tr *ProductsResourse) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product enteties.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := product.IsValidPrice()
	if err != nil {
		if errors.Is(err, enteties.ErrNegativePriceValue) {
			http.Error(w, "Invalid price, the price must be positive.", http.StatusBadRequest)
			return
		} else if errors.Is(err, enteties.ErrZeroPriceValue) {
			http.Error(w, "Invalid price, the price must be >0.", http.StatusBadRequest)
			return
		}
	}

	id, err := tr.S.CreateOneProductDb(product)
	if err != nil {
		log.Error().Msgf("Failed to create product in database: %v", err)
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
	log.Debug().Msg("Calling get all")

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
	log.Debug().Msgf("Requested IDs: %s", idsStr)

	ids, err := utils.ProcessIfIdsInt(idsStr)
	if err != nil {
		http.Error(w, "Wrong ID format", http.StatusBadRequest)
		return
	}

	idsStr = utils.ConvertIntSliceToString(ids)

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

	ok := utils.IsEmptyFullProduct(product)
	if ok {
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
	id := r.Context().Value(middleware.IdKey).(int)
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

func (tr *ProductsResourse) GetByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	products, err := tr.S.SearchProductByName(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No products found", http.StatusNotFound)
			return
		}
		http.Error(w, "No products found", http.StatusInternalServerError)
		return
	}
	if products == nil {
		http.Error(w, "No products found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Error().Msgf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (tr *ProductsResourse) GetCategorisedProducts(w http.ResponseWriter, r *http.Request) {
	category, ok := r.Context().Value(middleware.CategoryKey).(string)
	if !ok {
		http.Error(w, "Category not found in context", http.StatusInternalServerError)
		return
	}
	fmt.Println(category)
	products, err := tr.S.CategorisedProducts(category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No products found in this category", http.StatusNotFound)
			return
		}
		http.Error(w, "No products found in this category", http.StatusInternalServerError)
		return
	}
	if products == nil {
		http.Error(w, "No products found in this category", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Error().Msgf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
