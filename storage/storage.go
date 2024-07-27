package storage

import (
	"errors"
	"products/enteties"

	"sync"
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
	s.m.Lock()
	defer s.m.Unlock()

	s.lastId++
	p.ID = s.lastId
	s.allProducts[p.ID] = p
	return p.ID
}

func (s *Storage) GetAllProducts() {

}

func (s *Storage) DeleteProduct(ID int) bool {
	s.m.Lock()
	defer s.m.Unlock()

	if _, exists := s.allProducts[ID]; exists {
		delete(s.allProducts, ID)
		return true
	}
	return false
}

func (s *Storage) UpdateProduct() {

}

var (
	ErrProductNotFound = errors.New("product not found")
)

func (s *Storage) UpdateAvailability(id int, availability bool) error {
	product, exists := s.allProducts[id]
	if !exists {
		return ErrProductNotFound
	}

	product.IsAvailable = availability
	return nil
}
