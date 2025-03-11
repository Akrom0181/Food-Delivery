package entity

type UserLocation struct {
	Id         string  `json:"id"`
	UserId     string  `json:"user_id"`
	Address    string  `json:"address"`
	Entrance   uint    `json:"entrance"`
	Floor      int     `json:"floor"`
	DoorNumber uint    `json:"door_number"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type ListUserLocation struct {
	Items []UserLocation `json:"items"`
	Count int            `json:"count"`
}
