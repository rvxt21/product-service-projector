package resources

import (
	"encoding/json"
	"net/http"
	"product-service-projector/enteties"
	"product-service-projector/storage"
)

type ProductsResourse struct {
	S *storage.Storage
}

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

}

func (tr *ProductsResourse) UpdateProduct(w http.ResponseWriter, r *http.Request) {

}
