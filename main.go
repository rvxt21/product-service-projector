package main

import (
	"fmt"
	"log"
	"net/http"
	"products/resources"
	"products/storage"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Product Service Project!")

	connStr := "postgres://TemporaryMainuser:TemporaryPasw@database:5432/products?sslmode=disable"

	store, err := storage.New(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer store.db.Close()

	r := mux.NewRouter()
	productResource := &resources.ProductsResourse{S: store}

	productResource.RegisterRoutes(r)

	fmt.Println("Starting server at :8080")
	errServ := http.ListenAndServe(":8080", r)
	if errServ != nil {
		fmt.Println("Error happened", err.Error())
		return
	}
}
