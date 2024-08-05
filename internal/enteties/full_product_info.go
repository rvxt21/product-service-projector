package enteties

type FullProductInfo struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Price               int    `json:"price"`
	Quantity            int    `json:"quantity"`
	Category            int    `json:"category_id"`
	IsAvailable         bool   `json:"is_available"`
	CategoryName        string `json:"category_name"`
	CategoryDescription string `json:"category_description"`
}
