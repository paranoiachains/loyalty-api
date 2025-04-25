package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/order-service/internal/auth"
	"github.com/paranoiachains/loyalty-api/order-service/internal/database"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
)

func LoadOrder(c *gin.Context) {
	// get jwt token
	token, err := c.Cookie("jwt_token")
	if err != nil {
		logger.Log.Error("get cookie", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// retrieve userID from token
	userID := auth.GetUserID(token)
	if userID == -1 {
		logger.Log.Error("token is not valid")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

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

	err = database.DB.CreateAccrual(c.Request.Context(), accrualOrderID, userID)
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

	order := models.Accrual{
		UserID:         userID,
		AccrualOrderID: accrualOrderID,
	}

	data, err := json.Marshal(&order)
	if err != nil {
		logger.Log.Error("marshal accrual struct", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	messaging.OrderKafka.Send(data)
}

func GetOrder(c *gin.Context) {

}
