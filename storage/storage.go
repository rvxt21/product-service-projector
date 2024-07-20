package storage

import (
	"product-service-projector/enteties"
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

func (s *Storage) CreateOneProduct(p enteties.Product) int64 {
	s.m.Lock()
	defer s.m.Unlock()

	s.lastId++
	p.ID = s.lastId
	s.allProducts[p.ID] = p
	return p.ID
}

func (s *Storage) GetAllProducts() {

}

func (s *Storage) DeleteProduct() {

}

func (s *Storage) UpdateProduct() {

}
