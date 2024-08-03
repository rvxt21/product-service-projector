package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"products/resources"
	"products/storage"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Product Service Project!")

	connStr := os.Getenv("POSTGRES_CONN_STR")
	if connStr == "" {
		log.Fatal("Environment variable POSTGRES_CONN_STR is required")
	}

	store, err := storage.New(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer store.DB.Close()

	err = store.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	productResource := &resources.ProductsResourse{S: store}
	productResource.RegisterRoutes(r)

	fmt.Println("Starting server at :8080")
	errServ := http.ListenAndServe(":8080", r)
	if errServ != nil {
		fmt.Println("Error happened, %v", errServ.Error)
		return
	}
}
