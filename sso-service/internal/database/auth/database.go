package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
)

var (
	ErrUniqueUsername = errors.New("unique username must be set")
)

type Storage struct {
	db *sql.DB
}

func NewStorage(databaseDSN string) (*Storage, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s Storage) SaveUser(ctx context.Context, login string, passHash []byte) (uid int64, err error) {
	query := `
	INSERT INTO users(login, password)
	VALUES ($1, $2)
	`
	logger.Log.Info("saving user...")

	_, err = s.db.ExecContext(ctx, query, login, string(passHash))
	if err != nil {
		return 0, err
	}

	row := s.db.QueryRowContext(ctx, `SELECT user_id FROM users WHERE login = $1`)
	err = row.Scan(&uid)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, ErrUniqueUsername
		}
	}

	return uid, nil
}

func (s Storage) User(ctx context.Context, login string) (*models.User, error) {
	query := `
	SELECT user_id, login, password, balance, withdrawn
	FROM users
	WHERE login = $1
	`
	logger.Log.Info("retrieving user from db...", zap.String("user", login))

	row := s.db.QueryRowContext(ctx, query, login)

	var user models.User
	if err := row.Scan(&user.UserID, &user.Username, &user.Password, &user.Balance, &user.Withdrawn); err != nil {
		logger.Log.Error("retrieve user", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("user retrieved successfully!", zap.String("user", login))

	return &user, nil
}
