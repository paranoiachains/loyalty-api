package sso

import (
	"context"
	"errors"
	"fmt"

	sso_grpc "github.com/paranoiachains/loyalty-api/grpc-service/gen/go/sso"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ErrWrongPassword     = errors.New("wrong password")
	ErrUserAlreadyExists = errors.New("such user already exists")
)

type AuthClient struct {
	authClient sso_grpc.AuthClient
}

func New(address string) (*AuthClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := sso_grpc.NewAuthClient(conn)

	return &AuthClient{authClient: client}, nil
}

func (c *AuthClient) Login(ctx context.Context, login string, password string) (string, error) {
	resp, err := c.authClient.Login(ctx, &sso_grpc.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		logger.Log.Error("login", zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			logger.Log.Debug("convert err to status", zap.Bool("ok", ok))
			switch st.Code() {
			case codes.PermissionDenied:
				logger.Log.Debug("login (permission denied error)")
				return "", ErrWrongPassword
			default:
				return "", fmt.Errorf("unexpected grpc error: %w", err)
			}
		} else {
			return "", err
		}
	}

	return resp.Token, nil
}

func (c *AuthClient) RegisterNewUser(ctx context.Context, login string, password string) (int64, error) {
	resp, err := c.authClient.Register(ctx, &sso_grpc.RegisterRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return 0, ErrUserAlreadyExists
			default:
				return 0, fmt.Errorf("unexpected grpc error: %w", err)
			}
		} else {
			return 0, err
		}
	}

	return resp.UserId, nil
}
