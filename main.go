package main

import "fmt"

func main() {
	fmt.Println("Product Service Project!")

	store := storage.NewStorage()
	productResource := &resources.ProductsResource{Storage: store}

	http.HandleFunc("/products/create", productResource.CreateProduct)
	http.HandleFunc("/products/delete", productResource.DeleteProduct)

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080, nil)
}
