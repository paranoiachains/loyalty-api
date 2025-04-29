package withdraw

import (
	"context"
	"errors"
	"time"

	"github.com/paranoiachains/loyalty-api/grpc-service/gen/go/sso"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	database "github.com/paranoiachains/loyalty-api/sso-service/internal/database/withdraw"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	sso.UnimplementedWithdrawalsServer
	withdraw Withdraw
}

type Withdraw interface {
	TopUp(
		ctx context.Context,
		userID int64,
		sum float64,
	) error
	Balance(
		ctx context.Context,
		userID int64,
	) (current float64, withdrawn float64, err error)
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

func Register(gRPCServer *grpc.Server, withdraw Withdraw) {
	sso.RegisterWithdrawalsServer(gRPCServer, &serverAPI{withdraw: withdraw})
}

func (s *serverAPI) TopUp(
	ctx context.Context,
	in *sso.TopUpRequest,
) (*sso.TopUpResponse, error) {
	logger.Log.Info("balance top up (grpc level)", zap.Int64("user_id", in.UserId), zap.Float64("sum", in.Sum))

	if err := s.withdraw.TopUp(ctx, in.UserId, in.Sum); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso.TopUpResponse{}, nil
}

func (s *serverAPI) Balance(
	ctx context.Context,
	in *sso.BalanceRequest,
) (*sso.BalanceResponse, error) {
	logger.Log.Info("balance (grpc level)", zap.Int64("user_id", in.UserId))

	current, withdrawn, err := s.withdraw.Balance(ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso.BalanceResponse{Current: current, Withdrawn: withdrawn}, nil
}

func (s *serverAPI) Withdraw(
	ctx context.Context,
	in *sso.WithdrawRequest,
) (*sso.WithdrawResponse, error) {
	if err := s.withdraw.Withdraw(ctx, in.Order, in.UserId, in.Sum); err != nil {
		if errors.Is(err, database.ErrNotEnough) {
			return nil, status.Error(codes.Canceled, "not enough points")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso.WithdrawResponse{}, nil
}

func (s *serverAPI) Withdrawals(
	ctx context.Context,
	in *sso.WithdrawalsRequest,
) (*sso.WithdrawalsResponse, error) {
	withdrawals, err := s.withdraw.Withdrawals(ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	withdrawalPtrs := make([]*sso.Withdrawal, 0, len(withdrawals))

	for i := range withdrawals {
		withdrawal := &sso.Withdrawal{
			Order:       int64(withdrawals[i].OrderID),
			Sum:         withdrawals[i].Sum,
			ProcessedAt: withdrawals[i].ProcessedTime.Format(time.RFC3339),
		}

		withdrawalPtrs = append(withdrawalPtrs, withdrawal)
	}

	return &sso.WithdrawalsResponse{Withdrawals: withdrawalPtrs}, nil
}
