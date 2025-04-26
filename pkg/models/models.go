package models

import (
	"time"
)

type User struct {
	UserID    int     `json:"user_id"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	Balance   float64 `json:"balance"`
	Withdrawn float64 `json:"withdrawn"`
}

type Accrual struct {
	AccrualOrderID int       `json:"accrual_order_id"`
	UserID         int       `json:"user_id"`
	Status         string    `json:"status"`
	Accrual        float64   `json:"accrual"`
	UploadTime     time.Time `json:"uploaded_at"`
}

type AccrualStatusUpdate struct {
	OrderID int    `json:"accrual_order_id"`
	Status  string `json:"status"`
}

type Withdrawal struct {
	OrderID       int       `json:"order_id"`
	UserID        int       `json:"user_id"`
	Sum           float64   `json:"sum"`
	ProcessedTime time.Time `json:"processed_at"`
}
