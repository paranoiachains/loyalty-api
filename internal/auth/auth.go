package auth

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/paranoiachains/loyalty-api/internal/logger"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TokenExp = time.Hour * 3

var SecretKey = os.Getenv("SECRET_KEY")

// building jwt token with userID payload
func BuildJWTString(userID int) (string, error) {
	logger.Log.Info("building jwt token...")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		logger.Log.Error("signing token", zap.Error(err))
		return "", err
	}

	logger.Log.Info("successfully built jwt token!")
	return tokenString, nil
}

func GetUserID(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			return []byte(SecretKey), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}
	return claims.UserID
}

func SetCookies(c *gin.Context, userID int) error {
	logger.Log.Info("setting cookies for client...")
	token, err := BuildJWTString(userID)
	if err != nil {
		return err
	}
	c.SetCookie(
		"jwt_token",
		token,
		int(TokenExp),
		"/",
		"",
		false,
		true,
	)
	logger.Log.Info("cookies are set!")
	return nil
}
