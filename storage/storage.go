package storage

import (
	"sync"
)

type Storage struct {
	m           sync.Mutex
	lastId      int
	allProducts map[int]enteties.Task
}

func NewStorage() *Storage {
	return &Storage{
		allProducts: make(map[int]enteties.Product),
	}
}

func (s *Storage) CreateOneProduct() {

}

func (s *Storage) GetAllProducts() {

}

func (s *Storage) DeleteProduct() {

}

func (s *Storage) UpdateProduct() {

}

func (s *Storage) CreateOneProduct() {

}
