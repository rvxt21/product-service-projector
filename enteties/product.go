package enteties

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Category    string  `json:"category"`
	IsAvailable bool    `json:"is_available"`
}

type Catalogue struct {
	Products map[string]Product `json:"catalogue"`
}
