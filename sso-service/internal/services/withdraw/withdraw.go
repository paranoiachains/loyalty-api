package withdraw

import (
	"context"
	"errors"
	"fmt"

	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	database "github.com/paranoiachains/loyalty-api/sso-service/internal/database/withdraw"
	"go.uber.org/zap"
)

// Interfaces which must be implemented by Storage struct
type BalanceGetter interface {
	TopUp(
		ctx context.Context,
		userID int64,
		sum float64,
	) error
	Balance(
		ctx context.Context,
		userID int64,
	) (current float64, withdrawn float64, err error)
}

type Withdrawer interface {
	Withdraw(
		ctx context.Context,
		order int64,
		userID int64,
		sum float64,
	) error
	Withdrawals(
		ctx context.Context,
		userID int64,
	) ([]models.Withdrawal, error)
}

type Withdraw struct {
	balanceGetter BalanceGetter
	withdrawer    Withdrawer
}

func New(
	balanceGetter BalanceGetter, withdrawer Withdrawer,
) *Withdraw {
	return &Withdraw{
		balanceGetter: balanceGetter,
		withdrawer:    withdrawer,
	}
}

func (w *Withdraw) TopUp(
	ctx context.Context,
	userID int64,
	sum float64,
) error {
	logger.Log.Info("balance top up (service lvl)", zap.Int64("user_id", userID), zap.Float64("sum", sum))

	if err := w.balanceGetter.TopUp(ctx, userID, sum); err != nil {
		logger.Log.Error("top up", zap.Error(err))
		return err
	}

	return nil
}

func (w *Withdraw) Balance(
	ctx context.Context,
	userID int64,
) (current float64, withdrawn float64, err error) {
	logger.Log.Info("balance (service level)", zap.Int64("user_id", userID))

	current, withdrawn, err = w.balanceGetter.Balance(ctx, userID)
	if err != nil {
		logger.Log.Error("get balance", zap.Error(err))
		return 0, 0, err
	}

	logger.Log.Info("balance return (service level)", zap.Float64("current", current), zap.Float64("withdrawn", withdrawn))

	return current, withdrawn, nil
}

func (w *Withdraw) Withdraw(
	ctx context.Context,
	order int64,
	userID int64,
	sum float64,
) error {
	logger.Log.Info("withdrawing...", zap.Int64("order_id", order), zap.Int64("userID", userID), zap.Float64("sum", sum))

	if err := w.withdrawer.Withdraw(ctx, order, userID, sum); err != nil {
		logger.Log.Error("withraw", zap.Error(err))

		if errors.Is(err, database.ErrNotEnough) {
			return fmt.Errorf("%w", err)
		}
		return err
	}

	return nil
}

func (w *Withdraw) Withdrawals(
	ctx context.Context,
	userID int64,
) ([]models.Withdrawal, error) {
	withdrawals, err := w.withdrawer.Withdrawals(ctx, userID)
	logger.Log.Info("getting withdrawals...", zap.Int64("userID", userID))

	if err != nil {
		logger.Log.Info("get withdrawals...", zap.Error(err))
		return nil, err
	}

	return withdrawals, nil
}
