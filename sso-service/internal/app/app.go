package app

import (
	"time"

	"github.com/paranoiachains/loyalty-api/pkg/flags"
	grpcapp "github.com/paranoiachains/loyalty-api/sso-service/internal/app/grpc"
	databaseauth "github.com/paranoiachains/loyalty-api/sso-service/internal/database/auth"
	databasewithdraw "github.com/paranoiachains/loyalty-api/sso-service/internal/database/withdraw"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/services/auth"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/services/withdraw"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewAuth(grpcPort int, tokenTTL time.Duration) *App {
	db, err := databaseauth.NewStorage(flags.SSODatabaseDSN)
	if err != nil {
		panic(err)
	}

	authService := auth.New(db, db, tokenTTL)

	grpcApp := grpcapp.NewAuth(authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}

func NewWithdraw(grpcPort int) *App {
	db, err := databasewithdraw.NewStorage(flags.SSODatabaseDSN)
	if err != nil {
		panic(err)
	}

	withdrawService := withdraw.New(db, db)

	grpcApp := grpcapp.NewWithdraw(withdrawService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
