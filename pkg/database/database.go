package database

import (
	"context"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/paranoiachains/loyalty-api/pkg/models"
)

var (
	ErrUniqueUsername = errors.New("username already exists")
	ErrAlreadyExists  = errors.New("accrual for this order already exists for the same user")
	ErrAnotherUser    = errors.New("accrual for this order was already uploaded by other user")
)

type AccrualStorage interface {
	SetStatus(ctx context.Context, accrualOrderID int, status string) error
	UpdateAccrual(ctx context.Context, accrualOrderID int, accrual float64) error
	GetOrders(ctx context.Context, userID int) ([]models.Accrual, error)
	GetOrder(ctx context.Context, accrualOrderID int) (*models.Accrual, error)
	CreateAccrual(ctx context.Context, accrualOrderID int, userID int) (*models.Accrual, error)
}

type UserStorage interface {
	CreateUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type Storage interface {
	UserStorage
	AccrualStorage
}
