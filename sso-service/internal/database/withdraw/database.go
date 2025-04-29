package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
)

var (
	ErrNotEnough = errors.New("not enough points")
)

type Storage struct {
	db *sql.DB
}

func NewStorage(databaseDSN string) (*Storage, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		logger.Log.Error("db connect", zap.Error(err))
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s Storage) TopUp(
	ctx context.Context,
	userID int64,
	sum float64,
) error {
	query := `
	UPDATE users
	SET balance = balance + $1
	WHERE user_id = $2
	`
	logger.Log.Info("balance top up (db level)", zap.Int64("user_id", userID), zap.Float64("sum", sum))

	_, err := s.db.ExecContext(ctx, query, sum, userID)
	if err != nil {
		logger.Log.Error("top up db query", zap.Error(err))
		return err
	}

	logger.Log.Info("balance top up successful", zap.Int64("user_id", userID), zap.Float64("sum", sum))

	return nil
}

func (s Storage) Balance(
	ctx context.Context,
	userID int64,
) (current float64, withdrawn float64, err error) {
	query := `
	SELECT balance, withdrawn
	FROM users
	WHERE user_id = $1
	`
	logger.Log.Info("balance (db level)", zap.Int64("user_id", userID))

	row := s.db.QueryRowContext(ctx, query, userID)
	if err := row.Scan(&current, &withdrawn); err != nil {
		logger.Log.Error("scan rows (balance)", zap.Error(err))
		return 0, 0, err
	}

	logger.Log.Info("balance return (db level)", zap.Float64("current", current), zap.Float64("withdrawn", withdrawn))

	return current, withdrawn, nil
}

func (s Storage) Withdraw(
	ctx context.Context,
	order int64,
	userID int64,
	sum float64,
) error {
	queryWithdrawals := `
	INSERT INTO withdrawals(order_id, user_id, sum)
	VALUES ($1, $2, $3)
	`
	queryBalance := `
	UPDATE users
	SET balance = balance - $1, withdrawn = withdrawn + $1
	WHERE user_id = $2
	`

	logger.Log.Info("withdrawing, updating users table...")

	_, err := s.db.ExecContext(ctx, queryBalance, sum, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "balance_nonnegative" {
			return ErrNotEnough
		}
		logger.Log.Error("withdraw user balance", zap.Error(err))
		return err
	}

	_, err = s.db.ExecContext(ctx, queryWithdrawals, order, userID, sum)
	if err != nil {
		logger.Log.Error("add withdraw instance", zap.Error(err))
		return err
	}

	return nil
}

func (s Storage) Withdrawals(
	ctx context.Context,
	userID int64,
) ([]models.Withdrawal, error) {
	query := `
	SELECT order_id, sum, processed_at
	FROM withdrawals
	WHERE user_id = $1
	`
	logger.Log.Info("getting user withdrawals...", zap.Int64("user_ID", userID))

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Log.Error("retrieve withdrawals", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	withdrawals := make([]models.Withdrawal, 0)

	for rows.Next() {
		var withdrawal models.Withdrawal
		if err := rows.Scan(&withdrawal.OrderID, &withdrawal.Sum, &withdrawal.ProcessedTime); err != nil {
			logger.Log.Error("scan withdrawal", zap.Error(err))
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("rows iteration error", zap.Error(err))
		return nil, err
	}

	return withdrawals, nil
}
