package sso

import (
	"context"
	"errors"
	"fmt"

	sso_grpc "github.com/paranoiachains/loyalty-api/grpc-service/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ErrWrongPassword     = errors.New("wrong password")
	ErrUserAlreadyExists = errors.New("such user already exists")
)

type Client struct {
	authClient sso_grpc.AuthClient
}

func New(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := sso_grpc.NewAuthClient(conn)

	return &Client{authClient: client}, nil
}

func (c *Client) Login(ctx context.Context, login string, password string) (string, error) {
	resp, err := c.authClient.Login(ctx, &sso_grpc.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.PermissionDenied:
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

func (c *Client) RegisterNewUser(ctx context.Context, login string, password string) (int64, error) {
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
