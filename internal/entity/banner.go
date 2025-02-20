package entity

type Banner struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	ImageUrl  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BannerList struct {
	Items []Banner `json:"items"`
	Count int      `json:"count"`
}
