package auth

import (
	"context"
	"time"

	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/lib/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		login string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, login string) (*models.User, error)
}

type Auth struct {
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenTTL    time.Duration
}

func New(
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, login string, password string) (userID int64, err error) {
	logger.Log.Info("registering user...")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("generate hash from password", zap.Error(err))
		return 0, err
	}

	userID, err = a.usrSaver.SaveUser(ctx, login, passHash)
	if err != nil {
		logger.Log.Error("save user", zap.Error(err))
		return 0, err
	}

	return userID, nil
}

func (a *Auth) Login(ctx context.Context, login string, password string) (token string, err error) {
	logger.Log.Info("logging in", zap.String("login", login))

	user, err := a.usrProvider.User(ctx, login)
	if err != nil {
		return "", err
	}

	logger.Log.Debug("hash", zap.String("password", string(user.Password)))

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		logger.Log.Info("compare password bcrypt", zap.Error(err))
		return "", err
	}

	token, err = jwt.BuildJWTToken(user.UserID)
	if err != nil {
		return "", err
	}

	logger.Log.Info("user logged in", zap.String("user", login))

	return token, nil
}
