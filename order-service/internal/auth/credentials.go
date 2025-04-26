package auth

import (
	"context"

	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

// returns (nil, nil) if hash comparison went wrong
func (c Credentials) Authenticate(ctx context.Context, db database.Storage) (*models.User, error) {
	user, err := db.GetUserByUsername(ctx, c.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.Password)); err != nil {
		logger.Log.Warn("bcrypt compare", zap.Error(err))
		return nil, nil
	}
	return user, nil
}
