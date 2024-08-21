package resources

import (
	"net/http"
	"products/internal/middleware"
	"products/internal/storage"

	"github.com/gorilla/mux"
)

type ProductsResourse struct {
	S *storage.DBStorage
}

func (tr *ProductsResourse) RegisterRoutes(r *mux.Router) {
	r.Use(middleware.MockAuthenticationMiddleware)

	r.Handle("/products/{id}", middleware.AdminMiddleware(middleware.IdMiddleware(http.HandlerFunc(tr.DeleteProduct)))).Methods("DELETE")
	r.Handle("/products/{id}", middleware.AdminMiddleware(middleware.IdMiddleware(http.HandlerFunc(tr.UpdateProduct)))).Methods("PUT")
	r.Handle("/products/{id:[0-9]+}", middleware.IdMiddleware(middleware.AdminMiddleware(http.HandlerFunc(tr.GetProductByID)))).Methods("GET")
	r.Handle("/products/availability/{id}", middleware.AdminMiddleware(middleware.IdMiddleware(http.HandlerFunc(tr.UpdateAvailability)))).Methods("PATCH")
	r.Handle("/products", middleware.AdminMiddleware(http.HandlerFunc(tr.CreateProduct))).Methods("POST")

	r.HandleFunc("/products", tr.GetProducts).Methods("GET")
	r.HandleFunc("/products-by-ids", tr.GetProductsByIDS).Methods("GET")
	r.Handle("/products/{category:[a-zA-Z]+}", middleware.CategoryMiddleware(http.HandlerFunc(tr.GetCategorisedProducts))).Methods("GET")

	r.Handle("/categories", middleware.AdminMiddleware(http.HandlerFunc(tr.CreateCategory))).Methods("POST")
	r.HandleFunc("/categories", tr.GetAllCategories).Methods("GET")
	r.Handle("/categories/{idCategory}", middleware.IdMiddlewareCategory(http.HandlerFunc(tr.GetCategoryByID))).Methods("GET")
	r.Handle("/categories/{idCategory}", middleware.AdminMiddleware(middleware.IdMiddlewareCategory(http.HandlerFunc(tr.UpdateCategory)))).Methods("PUT")
	r.Handle("/categories/{idCategory}", middleware.AdminMiddleware(middleware.IdMiddlewareCategory(http.HandlerFunc(tr.DeleteCategory)))).Methods("DELETE")
	r.HandleFunc("/products/{id}", tr.GetProductByID).Methods("GET")
} //alternative for register routes
