package storage

import (
	"errors"
	"products/enteties"
	"sync"

	"github.com/rs/zerolog/log"
)

type Storage struct {
	m           sync.Mutex
	lastId      int
	allProducts map[int]enteties.Product
}

func NewStorage() *Storage {
	return &Storage{
		allProducts: make(map[int]enteties.Product),
	}
}

func (s *Storage) CreateOneProduct(p enteties.Product) int {
	const op = "storage.CreateProduct"
	s.m.Lock()
	defer s.m.Unlock()

	log.Info().Msgf("%s: creating product", op)
	s.lastId++
	p.ID = s.lastId
	s.allProducts[p.ID] = p
	return p.ID
}

// func (s *Storage) GetAllProducts() ([]enteties.Product, error) {

// }

func (s *Storage) DeleteProduct(ID int) bool {
	const op = "storage.DeleteProduct"
	s.m.Lock()
	defer s.m.Unlock()

	if _, exists := s.allProducts[ID]; exists {
		log.Info().Msgf("%s: deleting product %d", op, ID)
		delete(s.allProducts, ID)
		return true
	}
	return false
}

func (s *Storage) UpdateProduct(p enteties.Product) error {
	const op = "storage.UpdateProduct"
	s.m.Lock()
	defer s.m.Unlock()

	if _, exists := s.allProducts[p.ID]; !exists {
		log.Error().Msgf("%s: %s", op, ErrProductNotFound)
		return ErrProductNotFound
	}

	s.allProducts[p.ID] = p
	return nil
}

var (
	ErrProductNotFound = errors.New("product not found")
)

func (s *Storage) UpdateAvailability(id int, availability bool) error {
	const op = "storage.UpdateAvailability"
	product, exists := s.allProducts[id]
	if !exists {
		log.Error().Msgf("%s: %s", op, ErrProductNotFound)
		return ErrProductNotFound
	}

	product.IsAvailable = availability
	return nil
}
