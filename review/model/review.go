package model

import "time"

type Review struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	AccountID string `json:"account_id"`
	Rating    int       `json:"rating"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}