package storage

import (
	"database/sql"
	"fmt"
	"products/enteties"
	"sync"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type DBStorage struct {
	db *sql.DB
	m  sync.Mutex
}

func NewDBStorageDb(db *sql.DB) *DBStorage {
	return &DBStorage{db: db}
}

func New(connStr string) (*DBStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("openning database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	return &DBStorage{db: db}, nil
}

func (s *DBStorage) CreateOneProductDb(p enteties.Product) int {
	const op = "storage.CreateProduct"
	s.m.Lock()
	defer s.m.Unlock()

	log.Info().Msgf("%s: creating product", op)
	var id int
	err := s.db.QueryRow(
		"INSERT INTO products (name, description, price, quantity, category, is_available) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		p.Name, p.Description, p.Price, p.Quantity, p.Category, p.IsAvailable,
	).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to create product", op)
		return 0
	}

	return id
}

func (s *DBStorage) GetAllProductsDb() {
}

func (s *DBStorage) GetProductByIDDb(id int) {

}

func (s *DBStorage) DeleteProductDb(id int) (bool, error) {
	const op = "storage.DeleteProduct"

	query := `DELETE FROM products WHERE id=$1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		log.Error().Msgf("%s: deleting product: %v", op, err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error().Msgf("%s: checking rows affected: %v", op, err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return rowsAffected > 0, nil
}

func (s *Storage) UpdateProductBd(p enteties.Product) error {

}

func (db *DBStorage) UpdateProductAvailability(id int, availability bool) error {
	const op = "storage_db.UpdateProductAvailability"
	query := `UPDATE products SET availability = $1 WHERE ID = $2;`
	res, err := db.db.Exec(query, availability, id)
	if err != nil {
		log.Error().Err(err).Msgf("%s: %s", op, err)
		return err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		log.Info().Err(err).Msgf("%s: %s", op, err)
		return err
	}

	if rowAffected == 0 {
		log.Error().Msgf("%s: error to find product by ID", op)
		return sql.ErrNoRows
	}
	return nil
}
