package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type DBStorage struct {
	db *sql.DB
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
