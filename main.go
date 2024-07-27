package main

import (
	"fmt"
	"net/http"
	"products/middleware"
	"products/resources"
	"products/storage"
)

func main() {
	fmt.Println("Product Service Project!")
	mux := http.NewServeMux()
	store := storage.NewStorage()
	productResource := &resources.ProductsResourse{S: store}

	// productResource.RegisterRoutes(mux) //alternative for register routes

	http.HandleFunc("POST /products", productResource.CreateProduct)
	http.HandleFunc("DELETE /products", productResource.DeleteProduct)
	http.Handle("PATCH /products/{id}", middleware.IdMiddleware(http.HandlerFunc(productResource.UpdateAvailability)))

	fmt.Println("Starting server at :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error happened", err.Error())
		return
	}
}
