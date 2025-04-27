package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/paranoiachains/loyalty-api/order-service/internal/handlers/auth/models"
	"github.com/paranoiachains/loyalty-api/pkg/app"
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

		token, err := a.SSOClient.Login(context.Background(), creds.Login, creds.Password)
		if err != nil {
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
