package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"products/internal/resources"
	"products/internal/storage"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("Product Service Project!")
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	ctx := context.Background()

	res, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal().Msgf("Pinging redis: %v", err)
	}

	log.Info().Msgf("Pinged: %v", res)

	connStr := os.Getenv("POSTGRES_CONN_STR")
	if connStr == "" {
		log.Fatal().Msgf("Environment variable POSTGRES_CONN_STR is required")
	}

	store, err := storage.New(connStr, client)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	defer store.DB.Close()

	err = store.InitializeDB()
	if err != nil {
		log.Fatal().Msgf("%v", err)
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
