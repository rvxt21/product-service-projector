package storage

import (
	"database/sql"
	"errors"
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
	CREATE TABLE IF NOT EXISTS category (
		idCategory SERIAL PRIMARY KEY,
		nameCategory VARCHAR(100) NOT NULL,
		descriptionCategory TEXT NOT NULL
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

// нові функції // func (s *DBStorage) CreateCategory(category enteties.Category) (int, error) {}
// нові функції //func (s *DBStorage) UpdateCategory(category enteties.Category) error {}

func (s *DBStorage) GetAllProductsDb(limit, offset int) ([]enteties.Product, error) {
	const op = "storage.GetAllProducts"
	s.m.Lock()
	defer s.m.Unlock()

	query := `SELECT id, name, description, price, quantity, category, is_available 
			  FROM products 
			  LIMIT $1 
			  OFFSET $2`
	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to get all products", op)
		return nil, err
	}
	defer rows.Close()

	var products []enteties.Product
	for rows.Next() {
		var p enteties.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.Category, &p.IsAvailable); err != nil {
			log.Error().Err(err).Msgf("%s: unable to scan product", op)
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (s *DBStorage) GetProductByIDDb(id int) (enteties.Product, error) {
	const op = "storage.GetProductByIDDb"
	s.m.Lock()
	defer s.m.Unlock()

	var p enteties.Product
	err := s.DB.QueryRow(
		"SELECT id, name, description, price, quantity, category, is_available FROM products WHERE id = $1",
		id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.Category, &p.IsAvailable)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Msgf("%s: product not found with id %d", op, id)
			return enteties.Product{}, nil
		}
		log.Error().Err(err).Msgf("%s: unable to get product by id", op)
		return enteties.Product{}, err
	}

	return p, nil
}

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

var (
	ErrProductNotFound = errors.New("product not found")
)

func (s *DBStorage) UpdateProductAvailabilityDB(id int, availability bool) error {
	const op = "storage_db.UpdateProductAvailability"
	query := `UPDATE products SET is_available = $1 WHERE id = $2;`
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
