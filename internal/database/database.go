package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/paranoiachains/loyalty-api/internal/logger"
	"go.uber.org/zap"
)

// our database variable we will use throughout the process
var DB Storage

type Storage interface {
	Update() error
	Return() error
}

type PostgresStorage struct {
	*sql.DB
}

func (db PostgresStorage) Update() error {
	return nil
}

func (db PostgresStorage) Return() error {
	return nil
}

func Connect(databaseURI string) error {
	logger.Log.Info("connection to db", zap.String("database uri", databaseURI))
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		logger.Log.Error("init db connection error", zap.Error(err))
		return err
	}
	if err := db.Ping(); err != nil {
		logger.Log.Error("ping db", zap.Error(err))
		return err
	}
	logger.Log.Info("successufully connected to db")
	DB = PostgresStorage{db}
	return nil
}
