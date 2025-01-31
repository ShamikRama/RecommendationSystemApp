package user

import (
	"User_Service/internal/erorrs"
	"User_Service/internal/logger"
	"User_Service/internal/repository/model"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type DBUser struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewUserDB(db *sql.DB, logger *logger.Logger) *DBUser {
	return &DBUser{
		db:     db,
		logger: logger,
	}
}

func (r *DBUser) UpdateUser(ctx context.Context, user model.User) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("Error begin transaction", zap.String("operation", "sql.User.UpdateUser"), zap.Error(err))
		return 0, err
	}
	defer tx.Rollback()

	query := `
		UPDATE %s
		SET email = $1, password_hash = $2
		WHERE id = $3
		RETURNING id
	`
	query = fmt.Sprintf(query, model.UserTable)

	var id int

	err = tx.QueryRowContext(ctx, query, user.Email, user.Password, user.Id).Scan(&id)
	if err != nil {
		r.logger.Error("failed to query row", zap.String("operation", "sql.User.UpdateUser"), zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error("user not found", zap.String("operation", "sql.User.UpdateUser"), zap.Error(err))
			return 0, erorrs.ErrNoRows
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Error("Error commit the operation", zap.String("operation", "sql.User.UpdateUser"), zap.Error(err))
		return 0, err
	}

	return id, nil
}
