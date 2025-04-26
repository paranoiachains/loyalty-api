package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"go.uber.org/zap"
)

func GetOrder(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		accrualOrderID := c.Param("number")
		orderID, err := strconv.Atoi(accrualOrderID)
		if err != nil {
			logger.Log.Error("string to int", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		order, err := app.DB.GetOrder(context.Background(), orderID)
		if err != nil {
			logger.Log.Error("get order", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, order)
	}
}
