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

func Login(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds auth.Credentials

		if err := c.ShouldBindJSON(&creds); err != nil {
			logger.Log.Error("json request", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token, err := a.AuthClient.Login(context.Background(), creds.Login, creds.Password)
		if err != nil {
			if errors.Is(err, sso.ErrWrongPassword) {
				logger.Log.Error("login", zap.Error(err))
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			logger.Log.Error("login", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.SetCookie(
			"jwt_token",
			token,
			3600,
			"/",
			"",
			false,
			true,
		)

		c.String(http.StatusOK, "logged in successfully!")
	}
}
