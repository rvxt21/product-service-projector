package enteties

type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Quantity    int      `json:"quantity"`
	Category    Category `json:"category"`
	IsAvailable bool     `json:"is_available"`
}

type Category struct {
	IdCategory          int    `json:"idCategory"`
	NameCategory        string `json:"nameCategory"`
	DescriptionCategory string `json:"descriptionCategory"`
}
