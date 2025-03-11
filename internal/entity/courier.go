package entity

type Courier struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Status      string  `json:"status" example:"active, busy, inactive"` // active, busy, inactive
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	LastUpdated string  `json:"last_updated"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type OrderCourier struct {
	ID         string `json:"id"`
	OrderID    string `json:"order_id"`
	CourierID  string `json:"courier_id"`
	Status     string `json:"status" example:"pending, accepted, delivered"` // pending, accepted, delivered
	AssignedAt string `json:"assigned_at"`
}
