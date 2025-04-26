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

func SendStatus(status string, orderID int, p *LoyaltyProcessor) error {
	logger.Log.Info("sending status...", zap.String("status", status), zap.Int("orderID", orderID))
	statusMessage := models.AccrualStatusUpdate{
		OrderID: orderID,
		Status:  status,
	}

	payload, err := json.Marshal(&statusMessage)
	if err != nil {
		logger.Log.Error("marshal status message", zap.Error(err))
		return err

	}
	p.StatusBroker.Send(payload)

	logger.Log.Info("status sent!", zap.String("status", status), zap.Int("orderID", orderID))

	return nil
}

type LoyaltyProcessor struct {
	DB           database.Storage
	Broker       messaging.MessageBroker
	StatusBroker messaging.MessageBroker
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

		date := order.UploadTime

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

		SendStatus(Processing, createdOrder.AccrualOrderID, &p)

		// imitate evaluation
		time.Sleep(15 * time.Second)

		// evaluate accrual
		logger.Log.Info("evaluating accrual...")
		accrual := Evaluate()
		logger.Log.Info("accrual evaluated!")

		err = p.DB.UpdateAccrual(context.Background(), createdOrder.AccrualOrderID, accrual)
		if err != nil {
			logger.Log.Error("update accrual", zap.Error(err))
			continue
		}

		// set status to 'PROCESSED'
		err = p.DB.SetStatus(ctx, createdOrder.AccrualOrderID, Processed)
		if err != nil {
			logger.Log.Error("set status (db)", zap.Error(err))
			continue
		}

		SendStatus(Processed, createdOrder.AccrualOrderID, &p)

		// retrieve order from db
		processedOrder, err := p.DB.GetOrder(context.Background(), createdOrder.AccrualOrderID)
		if err != nil {
			logger.Log.Error("get order", zap.Error(err))
			continue
		}

		processedOrder.UploadTime = date

		// send back to kafka processed data
		processedData, err := json.Marshal(processedOrder)
		if err != nil {
			logger.Log.Error("marshal json", zap.Error(err))
			continue
		}

		p.Broker.Send(processedData)
	}
}
