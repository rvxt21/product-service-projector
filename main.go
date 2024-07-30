package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"products/middleware"
	"products/resources"
	"products/storage"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Product Service Project!")

	connStr := "postgres://TemporaryMainuser:TemporaryPasw@database:5432/products?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	store := storage.NewStorage()
	productResource := &resources.ProductsResourse{S: store}

	// productResource.RegisterRoutes(mux) //alternative for register routes

	mux.HandleFunc("POST /products", productResource.CreateProduct)
	mux.Handle("DELETE /products/{id}", middleware.IdMiddleware(http.HandlerFunc(productResource.DeleteProduct)))
	mux.Handle("PATCH /products/availability/{id}", middleware.IdMiddleware(http.HandlerFunc(productResource.UpdateAvailability)))
	mux.Handle("/products/{id}", middleware.IdMiddleware(http.HandlerFunc(productResource.UpdateProduct)))
	mux.HandleFunc("GET /products", productResource.GetAll)
	mux.Handle("GET /products/{id}", middleware.IdMiddleware(http.HandlerFunc(productResource.GetByID)))

	fmt.Println("Starting server at :8080")
	errServ := http.ListenAndServe(":8080", mux)
	if errServ != nil {
		fmt.Println("Error happened", err.Error())
		return
	}
}
