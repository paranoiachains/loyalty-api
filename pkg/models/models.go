package models

import (
	"time"
)

type User struct {
	UserID    int64   `json:"user_id"`
	Username  string  `json:"username"`
	Password  []byte  `json:"password"`
	Balance   float64 `json:"balance"`
	Withdrawn float64 `json:"withdrawn"`
}

type Accrual struct {
	AccrualOrderID int        `json:"order"`
	UserID         int        `json:"user_id,omitempty"`
	Status         string     `json:"status"`
	Accrual        float64    `json:"accrual"`
	UploadTime     *time.Time `json:"uploaded_at,omitempty"`
}

type AccrualStatusUpdate struct {
	OrderID int    `json:"order"`
	Status  string `json:"status"`
}

type Withdrawal struct {
	OrderID       int       `json:"order"`
	UserID        int       `json:"user_id,omitempty"`
	Sum           float64   `json:"sum"`
	ProcessedTime time.Time `json:"processed_at"`
}
