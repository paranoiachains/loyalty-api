package process

import (
	"context"
	"encoding/json"

	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
)

type OrderProcessor struct {
	DB     database.Storage
	Broker messaging.MessageBroker
}

func (p OrderProcessor) Process(ctx context.Context) {
	logger.Log.Info("processor started!")
	for data := range p.Broker.Receive() {
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
	}
}
