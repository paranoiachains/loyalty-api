package auth

import (
	"context"
	"errors"

	sso "github.com/paranoiachains/loyalty-api/grpc-service/gen/go/sso"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/database"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/services/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(
		ctx context.Context,
		login string,
		password string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		login string,
		password string,
	) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *sso.LoginRequest,
) (*sso.LoginResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.auth.Login(ctx, in.Login, in.Password)
	if err != nil {
		if errors.Is(err, auth.ErrWrongPassword) {
			logger.Log.Debug("login", zap.Error(err))
			return nil, status.Error(codes.PermissionDenied, "wrong password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &sso.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	id, err := s.auth.RegisterNewUser(ctx, in.Login, in.Password)
	if err != nil {
		if errors.Is(err, database.ErrUniqueUsername) {
			return nil, status.Error(codes.AlreadyExists, "such username already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &sso.RegisterResponse{UserId: id}, nil

}
