package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/order-service/internal/auth"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"go.uber.org/zap"
)

func Login(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds auth.Credentials
		if err := c.ShouldBindJSON(&creds); err != nil {
			logger.Log.Error("json request", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		user, err := creds.Authenticate(c.Request.Context(), app.DB)
		if user == nil && err == nil {
			logger.Log.Info("validation error, couldn't authenticate")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			logger.Log.Error("authenticate", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = auth.SetCookies(c, user.UserID)
		if err != nil {
			logger.Log.Error("setting cookies", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.String(http.StatusOK, "logged in")
	}
}
