package process

import (
	"context"
	"encoding/json"
	"time"

	"github.com/paranoiachains/loyalty-api/pkg/app"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
)

const (
	Processing = "PROCESSING"
	Processed  = "PROCESSED"
	Invalid    = "INVALID"
)

func Process(ctx context.Context, app *app.App) {
	for data := range app.Kafka.Receive() {
		var order models.Accrual
		err := json.Unmarshal(data, &order)
		if err != nil {
			logger.Log.Error("unmarshal data", zap.Error(err))
			continue
		}

		// set order status to 'PROCESSING'
		err = app.DB.SetStatus(ctx, order.AccrualOrderID, Processing)
		if err != nil {
			logger.Log.Error("set status (db)", zap.Error(err))
			continue
		}

		// imitate evaluation
		time.Sleep(20 * time.Second)

		// evaluate accrual
		order.Accrual = Evaluate()

		// set status to 'PROCESSED'
		err = app.DB.SetStatus(ctx, order.AccrualOrderID, Processed)
		if err != nil {
			logger.Log.Error("set status (db)", zap.Error(err))
			continue
		}

		// send back to kafka processed data
		processedData, err := json.Marshal(order)
		if err != nil {
			logger.Log.Error("marshal json", zap.Error(err))
			continue
		}

		app.Kafka.Send(processedData)
	}
}
