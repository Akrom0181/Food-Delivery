package entity

type Category struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CategorySingleRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryList struct {
	Items []Category `json:"items"`
	Count int        `json:"count"`
}
