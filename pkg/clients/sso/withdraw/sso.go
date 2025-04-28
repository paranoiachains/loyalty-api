package sso

import (
	"context"
	"time"

	sso_grpc "github.com/paranoiachains/loyalty-api/grpc-service/gen/go/sso"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WithdrawalsClient struct {
	withdrawalsClient sso_grpc.WithdrawalsClient
}

func New(address string) (*WithdrawalsClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := sso_grpc.NewWithdrawalsClient(conn)

	return &WithdrawalsClient{withdrawalsClient: client}, nil
}

func (w *WithdrawalsClient) Balance(ctx context.Context, userID int64) (current float64, withdrawn float64, err error) {
	logger.Log.Info("grpc calling... (balance)")

	resp, err := w.withdrawalsClient.Balance(ctx, &sso_grpc.BalanceRequest{
		UserId: userID,
	})
	if err != nil {
		logger.Log.Error("grpc call (balance)", zap.Error(err))
		return 0, 0, err
	}

	return resp.Current, resp.Withdrawn, nil
}

func (w *WithdrawalsClient) Withdraw(ctx context.Context, order int64, userID int64, sum float64) error {
	logger.Log.Info("withdrawing grpc call...")

	_, err := w.withdrawalsClient.Withdraw(ctx, &sso_grpc.WithdrawRequest{
		Order:  order,
		UserId: userID,
		Sum:    sum,
	})
	if err != nil {
		logger.Log.Error("withdraw grpc call", zap.Error(err))
		return err
	}

	return nil
}

func (w *WithdrawalsClient) Withdrawals(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	logger.Log.Info("withdrawals grpc call...")

	resp, err := w.withdrawalsClient.Withdrawals(ctx, &sso_grpc.WithdrawalsRequest{
		UserId: userID,
	})
	if err != nil {
		logger.Log.Error("withdrawals grpc call", zap.Error(err))
		return nil, err
	}

	withdrawals := make([]models.Withdrawal, 0)

	for _, ptrWithdrawal := range resp.Withdrawals {
		var withdrawal models.Withdrawal

		withdrawal.OrderID = int(ptrWithdrawal.Order)
		withdrawal.Sum = ptrWithdrawal.Sum
		withdrawal.ProcessedTime, err = time.Parse(time.RFC3339, ptrWithdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	return withdrawals, nil
}
