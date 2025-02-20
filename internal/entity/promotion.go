package entity

import "time"

type Promotion struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	DiscountPercentage int       `json:"discount_percentage"`
	StartedAt          time.Time `json:"started_at"`
	ExpiresAt          time.Time `json:"expires_at"`
	CreatedAt          string    `json:"created_at"`
}

type PromotionSingleRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PromotionGetList struct {
	Items []Promotion `json:"items"`
	Count int         `json:"count"`
}
