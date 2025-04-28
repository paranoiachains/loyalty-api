package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/paranoiachains/loyalty-api/order-service/internal/handlers/auth/models"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	sso "github.com/paranoiachains/loyalty-api/pkg/clients/sso/auth"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"go.uber.org/zap"
)

func Register(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds auth.Credentials

		if err := c.ShouldBindJSON(&creds); err != nil {
			logger.Log.Error("json request", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		_, err := a.AuthClient.RegisterNewUser(context.Background(), creds.Login, creds.Password)
		if err != nil {
			if errors.Is(err, sso.ErrUserAlreadyExists) {
				logger.Log.Error("register user", zap.Error(err))
				c.AbortWithStatus(http.StatusConflict)
				return
			}
			logger.Log.Error("register user", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.String(http.StatusOK, "user registered successfully!")
	}

}
