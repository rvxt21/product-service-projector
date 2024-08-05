package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"products/internal/enteties"
	"sync"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
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
    	category INT,
    	is_available BOOLEAN,
    	FOREIGN KEY (category) REFERENCES categories(idCategory) ON DELETE SET NULL);`

	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
    idCategory SERIAL PRIMARY KEY,
    nameCategory VARCHAR(100) NOT NULL,
    descriptionCategory TEXT NOT NULL);`

	s.m.Lock()
	defer s.m.Unlock()

	if _, err := s.DB.Exec(createCategoriesTable); err != nil {
		return fmt.Errorf("creating categories table: %w", err)
	}

	if _, err := s.DB.Exec(createProductsTable); err != nil {
		return fmt.Errorf("creating products table: %w", err)
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
	query := `INSERT INTO products (name, description, price, quantity, category, is_available) 
			  VALUES ($1, $2, $3, $4, $5, $6) 
		      RETURNING id`
	err := s.DB.QueryRow(query, p.Name, p.Description, p.Price, p.Quantity, p.Category, p.IsAvailable).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // Foreign key violation
				log.Error().Err(err).Msgf("%s: violates foreign key constraint", op)
				log.Info().Msgf(err.Error())
				return 0, err
			default:
				log.Error().Err(err).Msgf("%s: unable to create product", op)
				return 0, err
			}
		} else {
			log.Error().Err(err).Msgf("%s: unable to create product", op)
			return 0, err
		}
	}
	return id, nil
}

func (s *DBStorage) GetAllProductsDb(limit, offset int) ([]enteties.FullProductInfo, error) {
	const op = "storage.GetAllProducts"
	s.m.Lock()
	defer s.m.Unlock()

	query := `SELECT
					p.id AS product_id,
					p.name AS product_name,
					p.description AS product_description,
					p.price AS product_price,
					p.quantity AS product_quantity,
					p.is_available AS product_is_available,
					c.idCategory AS category_id,
					c.nameCategory AS category_name,
					c.descriptionCategory AS category_description
			FROM products p
			JOIN categories c ON p.category = c.idCategory
			LIMIT $1
			OFFSET $2;`
	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to get all products", op)
		return nil, err
	}
	defer rows.Close()

	var products []enteties.FullProductInfo
	for rows.Next() {
		var p enteties.FullProductInfo
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.IsAvailable, &p.Category, &p.CategoryName, &p.CategoryDescription); err != nil {
			log.Error().Err(err).Msgf("%s: unable to scan product", op)
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (s *DBStorage) GetProductsByIDSDB(selectingIds string) ([]enteties.FullProductInfo, error) {
	const op = "storage.GetProductsByIDS"

	s.m.Lock()
	defer s.m.Unlock()

	query := `SELECT
					p.id AS product_id,
					p.name AS product_name,
					p.description AS product_description,
					p.price AS product_price,
					p.quantity AS product_quantity,
					p.is_available AS product_is_available,
					c.idCategory AS category_id,
					c.nameCategory AS category_name,
					c.descriptionCategory AS category_description
			FROM products p
			JOIN categories c ON p.category = c.idCategory
			WHERE p.id IN (` + selectingIds + ");"

	rows, err := s.DB.Query(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to get range of products", op)
		return nil, err
	}
	defer rows.Close()

	var products []enteties.FullProductInfo
	for rows.Next() {
		var p enteties.FullProductInfo
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.IsAvailable, &p.Category, &p.CategoryName, &p.CategoryDescription); err != nil {
			log.Error().Err(err).Msgf("%s: unable to scan product", op)
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil

}

func (s *DBStorage) GetProductByIDDb(id int) (enteties.FullProductInfo, error) {
	const op = "storage.GetProductByIDDb"
	s.m.Lock()
	defer s.m.Unlock()

	var p enteties.FullProductInfo
	query := `SELECT
					p.id AS product_id,
					p.name AS product_name,
					p.description AS product_description,
					p.price AS product_price,
					p.quantity AS product_quantity,
					p.is_available AS product_is_available,
					c.idCategory AS category_id,
					c.nameCategory AS category_name,
					c.descriptionCategory AS category_description
			FROM products p
			JOIN categories c ON p.category = c.idCategory 
			WHERE id = $1`
	err := s.DB.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.IsAvailable, &p.Category, &p.CategoryName, &p.CategoryDescription)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Msgf("%s: product not found with id %d", op, id)
			return enteties.FullProductInfo{}, nil
		}
		log.Error().Err(err).Msgf("%s: unable to get product by id", op)
		return enteties.FullProductInfo{}, err
	}

	return p, nil
}

func (s *DBStorage) DeleteProductDb(id int) (bool, error) {
	const op = "storage.DeleteProduct"
	s.m.Lock()
	defer s.m.Unlock()

	query := `DELETE FROM products 
			  WHERE id=$1`

	result, err := s.DB.Exec(query, id)
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

func (s *DBStorage) UpdateProductBd(p enteties.Product) error {
	const op = "storage.UpdateProductBd"
	s.m.Lock()
	defer s.m.Unlock()

	query := `UPDATE products 
			  SET name = $1, description = $2, price = $3, quantity = $4, category = $5, is_available = $6 
			  WHERE id = $7`
	_, err := s.DB.Exec(query, p.Name, p.Description, p.Price, p.Quantity, p.Category, p.IsAvailable, p.ID)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to update product", op)
		return err
	}

	return nil
}

var (
	ErrProductNotFound = errors.New("product not found")
)

func (s *DBStorage) UpdateProductAvailabilityDB(id int, availability bool) error {
	const op = "storage_db.UpdateProductAvailability"
	query := `UPDATE products 
			  SET is_available = $1 
			  WHERE id = $2;`
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

func (s *DBStorage) GetAllCategoriesDb() ([]enteties.Category, error) {
	const op = "storage.GetAllCategories"
	s.m.Lock()
	defer s.m.Unlock()

	query := `SELECT * 
	          FROM categories`
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to get all categories", op)
		return nil, err
	}
	defer rows.Close()

	var categories []enteties.Category
	for rows.Next() {
		var c enteties.Category
		if err := rows.Scan(&c.IdCategory, &c.NameCategory, &c.DescriptionCategory); err != nil {
			log.Error().Err(err).Msgf("%s: unable to scan category", op)
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (s *DBStorage) CreateCategory(c enteties.Category) (int, error) {
	const op = "storage.CreateCategory"
	s.m.Lock()
	defer s.m.Unlock()

	query := `INSERT INTO categories (nameCategory, descriptionCategory) 
			  VALUES ($1, $2)
			  RETURNING idCategory`
	var id int
	err := s.DB.QueryRow(query, c.NameCategory, c.DescriptionCategory).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to create category", op)
		return 0, err
	}

	return id, nil
}

func (s *DBStorage) UpdateCategory(c enteties.Category) error {
	const op = "storage.UpdateCategory"
	s.m.Lock()
	defer s.m.Unlock()

	query := `UPDATE categories 
			  SET nameCategory=$1, descriptionCategory=$2 
			  WHERE idCategory=$3`
	_, err := s.DB.Exec(query, c.NameCategory, c.DescriptionCategory, c.IdCategory)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to update category", op)
		return err
	}

	return nil
}

func (s *DBStorage) GetCategoryByID(id int) (enteties.Category, bool, error) {
	const op = "storage.GetCategoryByID"
	s.m.Lock()
	defer s.m.Unlock()

	query := `SELECT idCategory, nameCategory, descriptionCategory 
			  FROM categories 
			  WHERE idCategory = $1`
	var c enteties.Category
	err := s.DB.QueryRow(query, id).Scan(&c.IdCategory, &c.NameCategory, &c.DescriptionCategory)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msgf("%s: unable to get category", op)
			return c, false, nil
		}
		return c, false, err
	}

	return c, true, nil
}

func (s *DBStorage) DeleteCategory(id int) (bool, error) {
	const op = "storage.DeleteCategory"
	s.m.Lock()
	defer s.m.Unlock()

	query := `DELETE FROM categories WHERE idCategory=$1`
	result, err := s.DB.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to delete category", op)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msgf("%s: unable to get rows affected", op)
		return false, err
	}

	return rowsAffected > 0, nil
}
