package storage

import (
	"database/sql"
	"errors"
	"products/internal/enteties"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

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
