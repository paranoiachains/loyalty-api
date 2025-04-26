package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type OrderStorage struct {
	*sql.DB
}

// sets an interface value to PostgresStorage
func Connect(databaseURI string) (OrderStorage, error) {
	logger.Log.Info("connecting to db...")
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		logger.Log.Error("init db connection error", zap.Error(err))
		return OrderStorage{}, err
	}
	if err := db.Ping(); err != nil {
		logger.Log.Error("ping db", zap.Error(err))
		return OrderStorage{}, err
	}
	logger.Log.Info("successufully connected to db")
	return OrderStorage{db}, nil
}

// creates user, returns user model
func (db OrderStorage) CreateUser(ctx context.Context, username, password string) (*models.User, error) {
	logger.Log.Info("creating user...")
	query := `
	INSERT INTO users (username, password, balance, withdrawn)
	VALUES ($1, $2, $3, $4);`

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
			return nil, database.ErrUniqueUsername
		}
		return nil, err
	}

	// also return user model
	var user models.User
	row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE username=$1;", username)
	err = row.Scan(&user.UserID, &user.Username, &user.Password, &user.Balance, &user.Withdrawn)
	if err != nil {
		return nil, err
	}

	logger.Log.Info("user successfully created!")
	return &user, nil
}

// return user model
func (db OrderStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	logger.Log.Info("retrieving user by username from db...")
	query := `SELECT user_id, username, password, balance, withdrawn FROM users WHERE username=$1;`

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

func (db OrderStorage) CreateAccrual(ctx context.Context, accrualOrderID int, userID int) (*models.Accrual, error) {
	logger.Log.Info("checking existing accrual...")

	var existingUserID int
	err := db.QueryRowContext(ctx, `
		SELECT user_id FROM accruals WHERE accrual_order_id = $1
	`, accrualOrderID).Scan(&existingUserID)

	if err == nil {
		if existingUserID == userID {
			return nil, database.ErrAlreadyExists
		}
		return nil, database.ErrAnotherUser
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	logger.Log.Info("creating accrual...")
	query := `
		INSERT INTO accruals (accrual_order_id, user_id, status)
		VALUES ($1, $2, $3);
	`

	_, err = db.ExecContext(ctx, query, accrualOrderID, userID, "NEW")
	if err != nil {
		return nil, err
	}
	logger.Log.Info("accrual created!")

	logger.Log.Info("returning accrual from db...")
	var order models.Accrual
	row := db.QueryRowContext(ctx,
		`SELECT accrual_order_id, user_id, status, accrual, uploaded_at FROM accruals WHERE accrual_order_id = $1`,
		accrualOrderID)
	err = row.Scan(&order.AccrualOrderID, &order.UserID, &order.Status, &order.Accrual, &order.UploadTime)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("accrual returned!")

	return &order, nil
}

func (db OrderStorage) SetStatus(ctx context.Context, accrualOrderID int, status string) error {
	logger.Log.Info("setting status...'",
		zap.Int("accrual_order_id", accrualOrderID),
		zap.String("status", status))

	query := `
	UPDATE accruals
	SET status = $1
	WHERE accrual_order_id = $2;
	`
	_, err := db.ExecContext(ctx, query, status, accrualOrderID)
	if err != nil {
		return err
	}

	logger.Log.Info("status set!",
		zap.Int("accrual_order_id", accrualOrderID),
		zap.String("status", status))

	return nil

}

func (db OrderStorage) UpdateAccrual(ctx context.Context, accrualOrderID int, accrual float64) error {
	query := `
	UPDATE accruals
	SET accrual = $1
	WHERE accrual_order_id = $2;
	`
	logger.Log.Info("updating accrual, starting tx...")
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Error("begin tx", zap.Error(err))
		return err
	}

	_, err = tx.ExecContext(ctx, query, accrual, accrualOrderID)
	if err != nil {
		logger.Log.Error("update accrual db query", zap.Error(err))
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error("commit tx", zap.Error(err))
		tx.Rollback()
		return err
	}

	logger.Log.Info("accrual updated!")
	return nil
}

func (db OrderStorage) GetOrders(ctx context.Context, userID int) ([]models.Accrual, error) {
	query := `
	SELECT accrual_order_id, user_id, status, accrual, uploaded_at
	FROM accruals
	WHERE user_id = $1
	ORDER BY uploaded_at ASC;
	`
	accruals := make([]models.Accrual, 0)

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Log.Error("query orders", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Accrual
		err := rows.Scan(
			&order.AccrualOrderID,
			&order.UserID,
			&order.Status,
			&order.Accrual,
			&order.UploadTime,
		)
		if err != nil {
			logger.Log.Error("scan rows", zap.Error(err))
		}
		accruals = append(accruals, order)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("rows iteration error", zap.Error(err))
		return nil, err
	}

	return accruals, nil
}

func (db OrderStorage) GetOrder(ctx context.Context, accrualOrderID int) (*models.Accrual, error) {
	return nil, nil
}
