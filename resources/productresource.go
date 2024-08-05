package resources

import (
	"net/http"
	"products/middleware"
	"products/storage"

	"github.com/gorilla/mux"
)

type ProductsResourse struct {
	S *storage.DBStorage
}

func (tr *ProductsResourse) RegisterRoutes(r *mux.Router) {
	r.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.DeleteProduct))).Methods("DELETE")
	r.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.UpdateProduct))).Methods("PUT")
	r.Handle("/products/availability/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.UpdateAvailability))).Methods("PATCH")
	r.HandleFunc("/products", tr.CreateProduct).Methods("POST")
	r.HandleFunc("/products", tr.GetAllProducts).Methods("GET")
	r.HandleFunc("/products-by-ids", tr.GetProductsByIDS).Methods("GET")
	r.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(tr.GetProductByID))).Methods("GET")

	r.HandleFunc("/categories", tr.CreateCategory).Methods("POST")
	r.HandleFunc("/categories", tr.GetAllCategories).Methods("GET")
	r.Handle("/categories/{idCategory}", middleware.IdMiddlewareCategory(http.HandlerFunc(tr.GetCategoryByID))).Methods("GET")
	r.Handle("/categories/{idCategory}", middleware.IdMiddlewareCategory(http.HandlerFunc(tr.UpdateCategory))).Methods("PUT")
	r.Handle("/categories/{idCategory}", middleware.IdMiddlewareCategory(http.HandlerFunc(tr.DeleteCategory))).Methods("DELETE")

} //alternative for register routes
