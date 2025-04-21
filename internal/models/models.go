package models

import "time"

type User struct {
	UserID    int     `json:"user_id"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	Balance   float64 `json:"balance"`
	Withdrawn float64 `json:"withdrawn"`
}

type Accural struct {
	AccuralOrderID int       `json:"accural_order_id"`
	UserID         int       `json:"user_id"`
	Status         string    `json:"status"`
	Accural        float64   `json:"accural"`
	UploadTime     time.Time `json:"uploaded_at"`
}

type Withdrawal struct {
	OrderID       int       `json:"order_id"`
	UserID        int       `json:"user_id"`
	Sum           float64   `json:"sum"`
	ProcessedTime time.Time `json:"processed_at"`
}
