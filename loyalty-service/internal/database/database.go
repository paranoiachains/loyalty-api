package database

import (
	"context"
	"database/sql"

	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
)

type LoyaltyStorage struct {
	*sql.DB
}

func Connect(databaseURI string) (LoyaltyStorage, error) {
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		logger.Log.Error("init db connection error", zap.Error(err))
		return LoyaltyStorage{}, err
	}
	if err := db.Ping(); err != nil {
		logger.Log.Error("ping db", zap.Error(err))
		return LoyaltyStorage{}, err
	}
	logger.Log.Info("successufully connected to db")
	return LoyaltyStorage{db}, nil
}

func (db LoyaltyStorage) CreateAccrual(ctx context.Context, accrualOrderID int, userID int) (*models.Accrual, error) {
	query := `
	INSERT INTO orders
	VALUES ($1, $2, $3, $4)
	`
	logger.Log.Info("creating order...")
	_, err := db.ExecContext(ctx, query, accrualOrderID, userID, "REGISTERED", 0)
	if err != nil {
		logger.Log.Error("create order (db)", zap.Error(err))
		return nil, err
	}
	logger.Log.Info("order created!")

	logger.Log.Info("returning order from db...")
	var order models.Accrual
	row := db.QueryRowContext(ctx,
		`SELECT order_id, user_id, status, accrual FROM orders WHERE order_id = $1`,
		accrualOrderID)
	err = row.Scan(&order.AccrualOrderID, &order.UserID, &order.Status, &order.Accrual)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("accrual returned!")
	return &order, nil
}

func (db LoyaltyStorage) SetStatus(ctx context.Context, accrualOrderID int, status string) error {
	query := `
	UPDATE orders
	SET status = $1
	WHERE order_id = $2 
	`
	logger.Log.Info("setting status...", zap.String("status", status))

	_, err := db.ExecContext(ctx, query, status, accrualOrderID)
	if err != nil {
		logger.Log.Error("set status (db)", zap.Error(err))
		return err
	}
	logger.Log.Info("status set!", zap.String("status", status))
	return nil
}

func (db LoyaltyStorage) UpdateAccrual(ctx context.Context, accrualOrderID int, accrual float64) error {
	query := `
	UPDATE orders
	SET accrual = $1
	WHERE order_id = $2;
	`
	logger.Log.Info("updating accrual, starting tx...")
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Error("begin tx", zap.Error(err))
		return err
	}

	_, err = tx.ExecContext(ctx, query, accrual, accrualOrderID)
	if err != nil {
		logger.Log.Error("update accrual db query", zap.Error(err))
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error("commit tx", zap.Error(err))
		tx.Rollback()
		return err
	}

	logger.Log.Info("accrual updated!")
	return nil
}

func (db LoyaltyStorage) GetOrders(ctx context.Context, userID int) ([]models.Accrual, error) {
	return nil, nil
}

func (db LoyaltyStorage) GetOrder(ctx context.Context, accrualOrderID int) (*models.Accrual, error) {
	query := `
	SELECT order_id, user_id, status, accrual
	FROM orders
	WHERE order_id = $1
	`

	var order models.Accrual
	row := db.QueryRowContext(ctx, query, accrualOrderID)
	err := row.Scan(&order.AccrualOrderID, &order.UserID, &order.Status, &order.Accrual)
	if err != nil {
		logger.Log.Error("get order", zap.Error(err))
		return nil, err
	}

	return &order, nil
}

func (db LoyaltyStorage) CreateUser(ctx context.Context, username, password string) (*models.User, error) {
	return nil, nil
}

func (db LoyaltyStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, nil
}
