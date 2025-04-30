package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"

	"github.com/paranoiachains/loyalty-api/pkg/app"
	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"go.uber.org/zap"
)

func LoadOrder(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, ok := c.Get("userID")
		if !ok {
			logger.Log.Error("get user id gin")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		userID := value.(int64)

		// read body, get order id
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Log.Error("read body", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := goluhn.Validate(string(body)); err != nil {
			logger.Log.Error("luhn validation", zap.Error(err))
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		// convert to int
		accrualOrderID, err := strconv.Atoi(string(body))
		if err != nil {
			logger.Log.Error("conv body to int", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		order, err := app.DB.CreateAccrual(ctx, accrualOrderID, int(userID))
		if err != nil {
			switch err {
			case database.ErrAlreadyExists:
				logger.Log.Warn("create accrual", zap.Error(err))
				c.String(http.StatusOK, "you've already loaded this order")
				return
			case database.ErrAnotherUser:
				logger.Log.Error("create accrual", zap.Error(err))
				c.String(http.StatusConflict, "this order was uploded by another user")
				return
			}
			logger.Log.Error("create accrual (unexpected)", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.String(http.StatusAccepted, "accrual instance created!")

		logger.Log.Info("marshalling order...")
		data, err := json.Marshal(&order)
		if err != nil {
			logger.Log.Error("marshal accrual struct", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		app.Kafka.Send(data)
	}
}

func GetOrders(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, ok := c.Get("userID")
		if !ok {
			logger.Log.Error("get user id gin")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		userID := value.(int64)

		orders, err := app.DB.GetOrders(context.Background(), int(userID))
		if err != nil {
			logger.Log.Error("get orders", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			c.String(http.StatusNoContent, "no orders")
		}

		c.JSON(http.StatusOK, orders)
	}
}
