package entity

type Notification struct {
	ID        string `json:"id"`
	OwnerId   string `json:"owner_id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	OwnerRole string `json:"ownerrole"`
	Message   string `json:"message"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type NotificationList struct {
	Notifications []Notification `json:"notifications"`
	Count         int            `json:"count"`
}
