package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"go.uber.org/zap"
)

func Balance(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, _ := c.Get("userID")
		userID := value.(int64)

		current, withdrawn, err := a.WithdrawClient.Balance(context.Background(), userID)
		if err != nil {
			logger.Log.Error("get balance", zap.Error(err))
			return
		}

		c.JSON(http.StatusOK, struct {
			current   float64
			withdrawn float64
		}{current: current, withdrawn: withdrawn})
	}
}

func Withdraw(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, _ := c.Get("userID")
		userID := value.(int64)

		withdrawal := struct {
			order int64
			sum   float64
		}{}
		c.ShouldBindJSON(&withdrawal)

		if err := goluhn.Validate(strconv.Itoa(int(withdrawal.order))); err != nil {
			logger.Log.Error("luhn not valid")
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if err := a.WithdrawClient.Withdraw(context.Background(), withdrawal.order, userID, withdrawal.sum); err != nil {
			logger.Log.Error("withdraw", zap.Error(err))
			return
		}

		c.String(http.StatusOK, "")

	}
}

func Withdrawals(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, _ := c.Get("userID")
		userID := value.(int64)

		withdrawals, err := a.WithdrawClient.Withdrawals(context.Background(), userID)
		if err != nil {
			logger.Log.Error("withdrawals", zap.Error(err))
			return
		}

		c.JSON(http.StatusOK, withdrawals)
	}
}
