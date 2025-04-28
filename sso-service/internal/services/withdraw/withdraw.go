package withdraw

import (
	"context"

	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
)

// Interfaces which must be implemented by Storage struct
type BalanceGetter interface {
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

func (w *Withdraw) Balance(
	ctx context.Context,
	userID int64,
) (current float64, withdrawn float64, err error) {
	logger.Log.Info("getting balance...", zap.Int64("userid", userID))

	current, withdrawn, err = w.balanceGetter.Balance(ctx, userID)
	if err != nil {
		logger.Log.Error("get balance", zap.Error(err))
		return 0, 0, err
	}

	return current, withdrawn, nil
}

func (w *Withdraw) Withdraw(
	ctx context.Context,
	order int64,
	userID int64,
	sum float64,
) error {
	logger.Log.Info("withdrawing...", zap.Int64("userID", order), zap.Float64("sum", sum))

	if err := w.withdrawer.Withdraw(ctx, order, userID, sum); err != nil {
		logger.Log.Error("withraw", zap.Error(err))
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
