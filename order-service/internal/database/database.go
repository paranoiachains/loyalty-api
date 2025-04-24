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
	"golang.org/x/crypto/bcrypt"
)

// if username already exists in db
var ErrUniqueUsername = errors.New("username already exists")

var ErrAlreadyExists = errors.New("accrual for this order already exists for the same user")
var ErrAnotherUser = errors.New("accrual for this order was already uploaded by other user")

// our database variable we will use throughout the process
var DB Storage

type Storage interface {
	CreateUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateAccrual(ctx context.Context, accrualOrderID int, userID int) error
}

type PostgresStorage struct {
	*sql.DB
}

// sets an interface value to PostgresStorage
func ConnectToPostgres(databaseURI string) error {
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

// creates user, returns user model
func (db PostgresStorage) CreateUser(ctx context.Context, username, password string) (*models.User, error) {
	logger.Log.Info("creating user...")
	query := `
	INSERT INTO users (username, password, balance, withdrawn)
	VALUES ($1, $2, $3, $4)`

	// bcrypt encryption of password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	// creating user
	_, err = db.ExecContext(ctx, query, username, encryptedPassword, 0, 0)
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
	row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE username=$1", username)
	err = row.Scan(&user.UserID, &user.Username, &user.Password, &user.Balance, &user.Withdrawn)
	if err != nil {
		return nil, err
	}

	logger.Log.Info("user successfully created!")
	return &user, nil
}

// return user model
func (db PostgresStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	logger.Log.Info("retrieving user by username from db...")
	query := `SELECT user_id, username, password, balance, withdrawn FROM users WHERE username=$1`

	var user models.User
	err := db.QueryRowContext(ctx, query, username).Scan(
		&user.UserID, &user.Username, &user.Password, &user.Balance, &user.Withdrawn,
	)
	if err != nil {
		return nil, err
	}

	logger.Log.Info("user successfully retrieved from db!")
	return &user, nil
}

func (db PostgresStorage) CreateAccrual(ctx context.Context, accrualOrderID int, userID int) error {
	logger.Log.Info("checking existing accrual...")

	var existingUserID int
	err := db.QueryRowContext(ctx, `
		SELECT user_id FROM accruals WHERE accrual_order_id = $1
	`, accrualOrderID).Scan(&existingUserID)

	if err == nil {
		if existingUserID == userID {
			return ErrAlreadyExists
		}
		return ErrAnotherUser
	}
	if err != sql.ErrNoRows {
		return err
	}

	logger.Log.Info("creating accrual...")
	query := `
		INSERT INTO accruals (accrual_order_id, user_id, status)
		VALUES ($1, $2, $3)
	`

	_, err = db.ExecContext(ctx, query, accrualOrderID, userID, "NEW")
	if err != nil {
		return err
	}

	return nil
}
