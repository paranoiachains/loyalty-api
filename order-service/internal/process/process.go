package process

import (
	"context"
	"encoding/json"

	ssowithdraw "github.com/paranoiachains/loyalty-api/pkg/clients/sso/withdraw"
	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
)

type OrderProcessor struct {
	DB             database.Storage
	Broker         messaging.MessageBroker
	StatusBroker   messaging.MessageBroker
	WithdrawClient *ssowithdraw.WithdrawalsClient
}

func (p OrderProcessor) Process(ctx context.Context) {
	logger.Log.Info("processor started!")

	brokerCh := p.Broker.Receive()
	statusCh := p.StatusBroker.Receive()

	for {
		select {
		case data, ok := <-brokerCh:
			if !ok {
				logger.Log.Warn("broker channel closed")
				return
			}
			var order models.Accrual
			logger.Log.Info("unmarshalling order...")
			err := json.Unmarshal(data, &order)
			if err != nil {
				logger.Log.Error("unmarshal order", zap.Error(err))
				continue
			}

			err = p.DB.UpdateAccrual(ctx, order.AccrualOrderID, order.Accrual)
			if err != nil {
				logger.Log.Error("update accrual", zap.Error(err))
				continue
			}

			logger.Log.Info("sending a top up request", zap.Int("user_id", order.UserID), zap.Float64("sum", order.Accrual))
			err = p.WithdrawClient.TopUp(ctx, int64(order.UserID), order.Accrual)
			if err != nil {
				logger.Log.Error("process top up call", zap.Error(err))
				continue
			}
		case data, ok := <-statusCh:
			if !ok {
				logger.Log.Warn("status broker channel closed")
				return
			}
			var statusUpdate models.AccrualStatusUpdate
			logger.Log.Info("unmarshalling status update...")
			err := json.Unmarshal(data, &statusUpdate)
			if err != nil {
				logger.Log.Error("unmarshal status update", zap.Error(err))
				continue
			}

			logger.Log.Info("status update received",
				zap.Int("order_id", statusUpdate.OrderID),
				zap.String("status", statusUpdate.Status),
			)

			err = p.DB.SetStatus(ctx, statusUpdate.OrderID, statusUpdate.Status)
			if err != nil {
				logger.Log.Error("set status", zap.Error(err))
				continue
			}
		}
	}
}
