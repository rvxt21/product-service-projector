package main

import (
	"fmt"
	"net/http"
	"product-service-projector/resources"
	"product-service-projector/storage"
)

func main() {
	fmt.Println("Product Service Project!")
	mux := http.NewServeMux()
	store := storage.NewStorage()
	productResource := &resources.ProductsResource{Storage: store}

	http.HandleFunc("/products/create", productResource.CreateProduct)
	http.HandleFunc("/products/delete", productResource.DeleteProduct)

	fmt.Println("Starting server at :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error happened", err.Error())
		return
	}
}
