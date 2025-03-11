package entity

type Order struct {
	ID             string       `json:"id"`
	UserID         string       `json:"user_id"`
	TotalPrice     float64      `json:"total_price"`
	Status         string       `json:"status" enums:"pending, confirmed, cancelled, preparing, picked_up, delivered" example:"pending"`
	DeliveryStatus string       `json:"delivery_status" enums:"olib ketish,yetkazib berish" example:"yetkazib berish"`
	Address        string       `json:"address"`
	Floor          int          `json:"floor"`
	DoorNumber     uint         `json:"door_number"`
	Entrance       uint         `json:"entrance"`
	Latitude       float64      `json:"latitude"`
	Longitude      float64      `json:"longitude"`
	OrderItems     []OrderItems `json:"order_items"`
	BranchId       string       `json:"branch_id"`
	CourierId      string       `json:"courier_id"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
}

type OrderList struct {
	Items []Order `json:"items"`
	Count int     `json:"count"`
}

type OrderItems struct {
	Id         string  `json:"id"`
	OrderId    string  `json:"order_id"`
	ProductId  string  `json:"product_id"`
	TotalPrice float64 `json:"total_price"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}
