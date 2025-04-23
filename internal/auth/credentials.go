package auth

import (
	"context"

	"github.com/paranoiachains/loyalty-api/internal/database"
	"github.com/paranoiachains/loyalty-api/internal/logger"
	"github.com/paranoiachains/loyalty-api/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

// returns nil, nil if hash comparison went wrong
func (c Credentials) Authenticate() (*models.User, error) {
	user, err := database.DB.GetUserByUsername(context.Background(), c.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.Password)); err != nil {
		logger.Log.Warn("bcrypt compare", zap.Error(err))
		return nil, nil
	}
	return user, nil
}
