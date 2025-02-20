package entity

type Product struct {
	Id          string  `json:"id"`
	CategoryId  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"image_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ProductSingleRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductList struct {
	Items []Product `json:"items"`
	Count int       `json:"count"`
}
