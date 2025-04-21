package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/internal/auth"
	"github.com/paranoiachains/loyalty-api/internal/database"
	"github.com/paranoiachains/loyalty-api/internal/logger"
	"go.uber.org/zap"
)

func Register(c *gin.Context) {
	if !strings.Contains(c.Request.Header.Get("Content-Type"), "application/json") {
		logger.Log.Info("content-type header",
			zap.String("want", "application/json"),
			zap.String("have", c.Request.Header.Get("Content-Type")))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		logger.Log.Error("read from body", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var creds auth.Credentials
	if err := json.Unmarshal(buf.Bytes(), &creds); err != nil {
		logger.Log.Error("unmarshalling body", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err := database.DB.CreateUser(context.Background(), creds)
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
