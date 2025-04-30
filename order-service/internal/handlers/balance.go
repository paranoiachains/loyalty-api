package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	sso "github.com/paranoiachains/loyalty-api/pkg/clients/sso/withdraw"
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
			Current   float64 `json:"current"`
			Withdrawn float64 `json:"withdrawn"`
		}{Current: current, Withdrawn: withdrawn})
	}
}

func Withdraw(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, _ := c.Get("userID")
		userID := value.(int64)

		withdrawal := struct {
			Order int64   `json:"order"`
			Sum   float64 `json:"sum"`
		}{}
		c.ShouldBindJSON(&withdrawal)

		logger.Log.Info("withdraw request handler lvl", zap.Int64("userID", userID), zap.Int64("order", withdrawal.Order), zap.Float64("sum", withdrawal.Sum))

		if err := goluhn.Validate(strconv.Itoa(int(withdrawal.Order))); err != nil {
			logger.Log.Error("luhn not valid")
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if err := a.WithdrawClient.Withdraw(context.Background(), withdrawal.Order, userID, withdrawal.Sum); err != nil {
			logger.Log.Error("withdraw", zap.Error(err))

			if errors.Is(err, sso.ErrNotEnough) {
				c.AbortWithStatus(http.StatusPaymentRequired)
				return
			}
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

		if len(withdrawals) == 0 {
			c.String(http.StatusNoContent, "no withdrawals")
		}

		c.JSON(http.StatusOK, withdrawals)
	}
}
