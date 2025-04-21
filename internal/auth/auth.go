package auth

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TOKEN_EXP = time.Hour * 3

var SECRET_KEY = os.Getenv("SECRET_KEY")

func BuildJWTString(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			return []byte(SECRET_KEY), nil
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
	token, err := BuildJWTString(userID)
	if err != nil {
		return err
	}
	c.SetCookie(
		"jwt_token",
		token,
		int(TOKEN_EXP),
		"/",
		"",
		false,
		true,
	)
	return nil
}
