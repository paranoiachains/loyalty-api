package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/paranoiachains/loyalty-api/internal/auth"
	"github.com/paranoiachains/loyalty-api/internal/logger"
	"github.com/paranoiachains/loyalty-api/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrUniqueUsername = errors.New("username already exists")

// our database variable we will use throughout the process
var DB Storage

type Storage interface {
	Update() error
	Return() error
	CreateUser(ctx context.Context, creds auth.Credentials) (*models.User, error)
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
	logger.Log.Info("connecting to db...")
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

func (db PostgresStorage) CreateUser(ctx context.Context, creds auth.Credentials) (*models.User, error) {
	logger.Log.Info("creating user...")
	query := `
	INSERT INTO users (username, password, balance, withdrawn)
	VALUES ($1, $2, $3, $4)`

	// bcrypt encryption of password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// creating user
	_, err = tx.ExecContext(ctx, query, creds.Username, encryptedPassword, 0, 0)
	if err != nil {
		var pgErr *pgconn.PgError
		// if username already taken
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUniqueUsername
		}
		return nil, err
	}

	// also return user model
	var user models.User
	row := tx.QueryRowContext(ctx, "SELECT * FROM users WHERE username=$1", creds.Username)
	err = row.Scan(&user.UserID, &user.Username, &user.Password, &user.Balance, &user.Withdrawn)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	logger.Log.Info("user successfully created!")
	return &user, nil
}
