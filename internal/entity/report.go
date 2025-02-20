package entity

type Report struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ReportList struct {
	Items []Report `json:"items"`
	Count int      `json:"count"`
}
