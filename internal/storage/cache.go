package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"products/internal/enteties"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

func (db *DBStorage) CacheProduct(product enteties.FullProductInfo) error {
	productData, err := json.Marshal(product)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("product:%d", product.ID)
	err = db.cache.Set(ctx, key, productData, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

var ErrNoProductInCache = errors.New("no product in cache")

func (db *DBStorage) GetCachedProduct(productId int) (*enteties.FullProductInfo, error) {
	key := fmt.Sprintf("product:%d", productId)
	result, err := db.cache.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Info().Msgf("Product with ID %d not found in cache", productId)
			return nil, ErrNoProductInCache
		}
		return nil, err
	}

	var product enteties.FullProductInfo
	if err := json.Unmarshal([]byte(result), &product); err != nil {
		return nil, err
	}
	return &product, nil
}
