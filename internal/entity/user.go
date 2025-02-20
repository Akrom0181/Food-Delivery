package entity

type User struct {
	ID          string `json:"id"`
	FullName    string `json:"full_name"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	UserType    string `json:"user_type"`
	UserRole    string `json:"user_role"`
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
	ProfilePic  string `json:"profile_picture"`
	Gender      string `json:"gender"`
	Bio         string `json:"bio"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserSingleRequest struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type UserList struct {
	Items []User `json:"users"`
	Count int    `json:"count"`
}
