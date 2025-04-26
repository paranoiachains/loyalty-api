package process

import (
	"context"
	"encoding/json"
	"time"

	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
)

const (
	Processing = "PROCESSING"
	Processed  = "PROCESSED"
	Invalid    = "INVALID"
)

type LoyaltyProcessor struct {
	DB     database.Storage
	Broker messaging.MessageBroker
}

func (p LoyaltyProcessor) Process(ctx context.Context) {
	logger.Log.Info("processor started!")
	for data := range p.Broker.Receive() {
		var order models.Accrual
		logger.Log.Info("unmarshalling order...")
		err := json.Unmarshal(data, &order)
		if err != nil {
			logger.Log.Error("unmarshal data", zap.Error(err))
			continue
		}

		createdOrder, err := p.DB.CreateAccrual(ctx, order.AccrualOrderID, order.UserID)
		if err != nil {
			logger.Log.Error("create order", zap.Error(err))
			continue
		}

		logger.Log.Info("order created", zap.String("status", createdOrder.Status))

		// set order status to 'PROCESSING'
		err = p.DB.SetStatus(ctx, order.AccrualOrderID, Processing)
		if err != nil {
			logger.Log.Error("set status (db)", zap.Error(err))
			continue
		}

		// imitate evaluation
		time.Sleep(15 * time.Second)

		// evaluate accrual
		logger.Log.Info("evaluating accrual...")
		order.Accrual = Evaluate()
		logger.Log.Info("accrual evaluated!")

		// set status to 'PROCESSED'
		err = p.DB.SetStatus(ctx, order.AccrualOrderID, Processed)
		if err != nil {
			logger.Log.Error("set status (db)", zap.Error(err))
			continue
		}

		// retrieve order from db
		processedOrder, err := p.DB.GetOrder(context.Background(), order.AccrualOrderID)
		if err != nil {
			logger.Log.Error("get order", zap.Error(err))
			continue
		}

		// send back to kafka processed data
		processedData, err := json.Marshal(processedOrder)
		if err != nil {
			logger.Log.Error("marshal json", zap.Error(err))
			continue
		}

		p.Broker.Send(processedData)
	}
}
