package app

import (
	"time"

	grpcapp "github.com/paranoiachains/loyalty-api/sso-service/internal/app/grpc"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/database"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(grpcPort int, tokenTTL time.Duration) *App {
	db, err := database.NewStorage("postgresql://postgres:postgres@postgres/sso_service?sslmode=disable")
	if err != nil {
		panic(err)
	}

	authService := auth.New(db, db, tokenTTL)

	grpcApp := grpcapp.New(authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
