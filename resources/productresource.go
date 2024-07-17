package resources

import "net/http"

type ProductsResourse struct {
	S *storage.Storage
}

func (tr *ProductsResourse) CreateTask(w http.ResponseWriter, r *http.Request) {

}

func (tr *ProductsResourse) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (tr *ProductsResourse) DeleteProduct(w http.ResponseWriter, r *http.Request) {

}

func (tr *ProductsResourse) UpdateProduct(w http.ResponseWriter, r *http.Request) {

}
