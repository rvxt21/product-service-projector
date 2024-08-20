package enteties

import "errors"

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	Category    int    `json:"category"`
	IsAvailable bool   `json:"is_available"`
}

var (
	ErrNegativePriceValue = errors.New("negative price")
	ErrZeroPriceValue     = errors.New("zero price value")
)

func (p Product) IsValidPrice() error {
	switch {
	case p.Price < 0:
		return ErrNegativePriceValue
	case p.Price == 0:
		return ErrZeroPriceValue
	default:
		return nil
	}
}
