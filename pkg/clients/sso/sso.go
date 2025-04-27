package sso

import (
	"context"

	sso_grpc "github.com/paranoiachains/loyalty-api/grpc-service/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		return "", err
	}

	return resp.Token, nil
}

func (c *Client) RegisterNewUser(ctx context.Context, login string, password string) (int64, error) {
	resp, err := c.authClient.Register(ctx, &sso_grpc.RegisterRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return 0, err
	}

	return resp.UserId, nil
}
