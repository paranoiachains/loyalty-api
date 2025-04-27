package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"go.uber.org/zap"
)

func BuildJWTToken(userID int64) (string, error) {
	logger.Log.Info("building jwt token...")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		logger.Log.Error("signing token", zap.Error(err))
		return "", err
	}

	logger.Log.Info("successfully built jwt token!")

	return tokenString, nil
}
