package enteties

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	Category    int    `json:"category"`
	IsAvailable bool   `json:"is_available"`
}
