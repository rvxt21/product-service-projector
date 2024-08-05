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
