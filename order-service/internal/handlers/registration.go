package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/order-service/internal/auth"
	"github.com/paranoiachains/loyalty-api/order-service/internal/database"
	"github.com/paranoiachains/loyalty-api/order-service/internal/logger"
	"go.uber.org/zap"
)

func Register(c *gin.Context) {
	var creds auth.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		logger.Log.Error("json request", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := database.DB.CreateUser(context.Background(), creds.Password, creds.Username)
	if err != nil {
		logger.Log.Error("creating user", zap.Error(err))

		if errors.Is(err, database.ErrUniqueUsername) {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = auth.SetCookies(c, user.UserID)
	if err != nil {
		logger.Log.Error("setting cookies", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, "user created")
}
