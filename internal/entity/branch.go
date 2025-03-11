package entity

type Branch struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BranchSingleRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BranchList struct {
	Items []Branch `json:"items"`
	Count int      `json:"count"`
}
