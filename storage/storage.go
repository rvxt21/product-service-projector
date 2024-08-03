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
	DB *sql.DB
	m  sync.Mutex
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{DB: db}
}

func New(connStr string) (*DBStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("openning database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	return &DBStorage{DB: db}, nil
}

func (s *DBStorage) InitializeDB() error {
	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		description TEXT,
		price INT,
		quantity INT,
		category VARCHAR(100),
		is_available BOOLEAN
	);`

	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT NOT NULL
	);`

	s.m.Lock()
	defer s.m.Unlock()

	if _, err := s.DB.Exec(createProductsTable); err != nil {
		return fmt.Errorf("creating products table: %w", err)
	}

	if _, err := s.DB.Exec(createCategoriesTable); err != nil {
		return fmt.Errorf("creating categories table: %w", err)
	}

	log.Info().Msg("Database initialized successfully")
	return nil
}

func (s *DBStorage) CreateOneProductDb(p enteties.Product) (int, error) {
	const op = "storage.CreateProduct"
	s.m.Lock()
	defer s.m.Unlock()

	log.Info().Msgf("%s: creating product", op)
	var id int
	err := s.DB.QueryRow(
		"INSERT INTO products (name, description, price, quantity, category, is_available) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		p.Name, p.Description, p.Price, p.Quantity, p.Category, p.IsAvailable,
	).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to create product", op)
		return 0, err
	}

	return id, nil
}

//func (s *DBStorage) GetAllProductsDb() {}

//func (s *DBStorage) GetProductByIDDb(id int) {}

func (s *DBStorage) DeleteProductDb(id int) (bool, error) {
	const op = "storage.DeleteProduct"
	s.m.Lock()
	defer s.m.Unlock()

	result, err := s.DB.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to delete product", op)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to get rows affected", op)
		return false, err
	}

	return rowsAffected > 0, nil
}

//func (s *Storage) UpdateProductBd(p enteties.Product) error {}

func (s *DBStorage) UpdateProductAvailability(id int, availability bool) error {
	const op = "storage_db.UpdateProductAvailability"
	query := `UPDATE products SET availability = $1 WHERE ID = $2;`
	res, err := s.DB.Exec(query, availability, id)
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
